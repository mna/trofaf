---
Date: 2012-04-13
Title: Propriétés calculées avec Backbone
Author: Martin Angers
Category: technologie
Description: La beauté de l'univers du code libre est que lorsqu'il manque une fonctionnalité, on peut se retrousser les manches, ouvrir le code et l'ajouter. Voici ma petite histoire des propriétés calculées avec Backbone.
---

Je travaille actuellement à temps plein sur mon application-en-devenir de [Compte Rendu Web][crw]. Cette immersion totale en Javascript, [node.js][node] et [Backbone][bb] est très intéressante. La beauté de l'univers du code libre dans lequel baignent ces technologies est que lorsqu'il manque une fonctionnalité, on peut se retrousser les manches, ouvrir le code et l'ajouter. Voici ma petite histoire des propriétés calculées (*computed properties*) avec Backbone.

### Pas encore une controverse Backbone vs Knockout...

Backbone et [Knockout][ko] sont deux superbes librairies (parmi tant d'autres!) de type [MV*][mvx] qui offrent une meilleure façon d'organiser et de découper le code client. Personnellement je préfère Backbone pour son côté minimaliste, *juste assez*, sans magie et très extensible. Il y a toutefois une fonctionnalité intéressante de Knockout qui m'a manqué récemment, soit la possibilité de créer des *computed observables*, dans le jargon de KO.

L'idée est simple, ayant une propriété "nom" et "prénom", on peut créer une propriété calculée "nomComplet" qui concatène les deux. Ça va plus loin, on peut même assigner une valeur à cette propriété calculée, et les propriétés "réelles" sous-jacentes peuvent être ainsi renseignées. J'avais un besoin pour ce type de comportement. Bien sûr, avec Backbone, il y a d'autres façons de s'en sortir, avec du code dans la vue (`Backbone.View`) qui s'abonne aux événements `change:nomAttribut` et qui réagit en conséquence, mais voilà, j'utilise aussi [Backbone.ModelBinding][bmb] pour lier les champs HTML aux attributs du modèle, et ça ne me disait pas du tout d'avoir une partie du modèle liée et une autre gérée manuellement dans la vue.

J'avais donc ce qui commençait à ressembler à des spécifications:

*	Une propriété calculée peut être basée sur un ou plusieurs attributs "réels" d'un modèle.
*	Lorsqu'un de ces attributs lance un événement "change", la propriété calculée lance son propre événement "change", afin de permettre aux observateurs de rafraîchir la donnée calculée.
*	La propriété calculée doit supporter la lecture et l'écriture.
*	La propriété calculée ne doit **pas** polluer le tableau des attributs "réels" du modèle (`Backbone.Model.attributes`), ni - par le fait même - le JSON envoyé au serveur lors des sauvegardes (ce sont des propriétés *calculées*, donc déduites à partir d'autres attributs, pas d'intérêt à les sauvegarder).
*	L'implémentation doit être compatible avec Backbone.ModelBinding.

### L'implémentation

Certains voudront peut-être [sauter au code][gist] pour revenir aux explications par la suite. Faites ça vite, je vous attends.

Sur le modèle comme tel, une seule nouvelle propriété est requise, soit `Model.computedProperties`, qui est une instance de `ComputedProperties`. Cette classe supporte trois méthodes: `add()`, `remove()` et `clear()`, et conserve les propriétés calculées dans un tableau `ComputedProperties.properties`. Rien de bien captivant à ce niveau, c'est plutôt sur la classe `ComputedProperty` (notez la nuance) que le tout se déroule.

Une propriété calculée est créée en passant 4 informations à `ComputedProperties.add()` (soit en un seul objet, soit en 4 paramètres distincts):

1.	`name`: le nom de la propriété calculée.
2.	`attributes`: un tableau (*array*) des attributs du modèle sur lesquels est basée la propriété calculée. La `ComputedProperty` s'abonnera aux événements "change:<attribut>" de ceux-ci pour déclencher son propre événement "change:<compProp>".
3.	`getter`: la fonction appelée pour obtenir la valeur de la propriété calculée. Celle-ci ne reçoit aucun paramètre et s'exécute dans le contexte du modèle (`this` est le modèle, permettant ainsi d'appeler `this.get("attr")` pour obtenir les valeurs des attributs).
4.	`setter`: la fonction appelée pour assigner la valeur à la propriété calculée, exécutée dans le contexte du modèle. Règle générale, celle-ci devrait assigner une valeur aux attributs de base, déclenchant ainsi le "change" de l'attribut qui déclenche le "change" de la propriété calculée (vous me suivez?). Si - étrangement - le *setter* n'assigne de valeur à aucun attribut relié à la propriété calculée, alors celui-ci devra déclencher manuellement le "change" de la propriété. La fonction reçoit deux paramètres, la nouvelle valeur et l'objet des options (le même que celui reçu par `Backbone.Model.set()`). Information importante à ce sujet: si les options spécifient `unset: true` (ou n'importe quelle autre option pertinente, tant qu'à ça), il est de la responsabilité du *setter* d'agir en conséquence.

Seul le `name` est obligatoire. La classe `ComputedProperty` offre d'autres méthodes utilitaires, telles que `isReadOnly()`, `isReadWrite()` et `isWriteOnly()`.

Reste maintenant à attacher le tout. **Backbone.ModelBinding** écoute les événements "change:attr" pour assigner la valeur aux éléments de la page Web (s'attendant à retrouver la valeur comme 2ème paramètre). L'implémentation des propriétés calculées répond à ce cas d'utilisation. Inversement, lorsque la valeur est modifiée dans le champ de la page Web, il utilise `model.set("attribut")`, alors qu'on veut justement éviter que la propriété calculée aboutisse dans les attributs. Il faut faire quelque chose pour gérer ce cas. Et ce serait bien aussi de pouvoir faire `model.get("nomPropCalculee")` et obtenir la valeur de celle-ci.

C'est pourquoi le modèle de base hérité de `Backbone.Model` substitue les méthodes `get()` et `set()` pour intercepter les cas où un nom de propriété calculée est reçu, et éviter qu'elle se retrouve dans les attributs et soit transférée au serveur.

J'ai écrit une cinquantaine de tests unitaires (avec [Mocha][] comme engin de test, [expect.js][expect] pour les assertions et [sinon.js][sinon] pour les [doubles de test][double], c'est une superbe combinaison que je recommande!), ça semble assez solide quoique rapidement codé et testé pour le moment. Il y a seulement quelques trucs à savoir et sur lesquels méditer:

*	D'abord, l'événement "change" de la propriété calculée est déclenché chaque fois qu'un attribut utilisé par cette propriété est modifié. Si plusieurs attributs sont modifiés en un seul appel à `Model.set()`, plusieurs "change" de la propriété calculée seront déclenchés.
*	Ensuite, il y a peut-être des cas plus complexes d'utilisation de `Backbone.Model.set()` ou `Backbone.Model.save()` qui m'ont échappé et qui fonctionnent mal avec la substitution du `set()`. Si vous en trouvez, m'en faire part, je maintiendrai le *gist* à jour. Je doute que le `get()` pose problème.
*	Le fait que les propriétés calculées ne se retrouvent pas dans les attributs était pour moi un objectif important, mais je peux imaginer certains cas où on voudrait à tout le moins les avoir dans le JSON (par exemple si on utilise `Backbone.Model.toJSON()` pour passer le résultat à un modèle de vue - un *template*). Ce serait très simple d'ajouter un paramètre dans une substitution de `toJSON()` pour indiquer que l'on désire y retrouver aussi ces propriétés, et lors de l'envoi au serveur par un `save()` interne à Backbone, ce paramètre ne serait pas fourni, donc les propriétés ne seraient pas envoyées.
*	Et finalement, il y a l'aspect conceptuel de la chose, qui éloigne le modèle de la représentation des données et l'approche de la représentation "humaine". Ça rend Backbone un peu plus MVVM que MVC, si on veut.

### Un gist pour le moment

J'ai mis le code source dans ce [gist][gist], car le tout fait à peine 100 [sloc][], il est probablement préférable de le copier et l'insérer dans son modèle Backbone de base (et de le minifier avec son code maison) plutôt que de le distribuer comme une librairie indépendante. Et en plus, le code assume qu'il hérite directement de Backbone.Model, alors que ce n'est pas forcément le cas. Qui sait, peut-être qu'un auteur de librairie Backbone de plus haut niveau sera intéressé à l'intégrer dans son cadre de développement (*framework*).

[crw]: http://www.compterenduweb.com/
[node]: http://nodejs.org/
[bb]: http://backbonejs.org/
[ko]: http://knockoutjs.com/
[mvx]: http://www.codeproject.com/Articles/42830/Model-View-Controller-Model-View-Presenter-and-Mod#_articleTop
[gist]: https://gist.github.com/2371954
[sloc]: http://en.wikipedia.org/wiki/Source_lines_of_code
[bmb]: https://github.com/derickbailey/backbone.modelbinding
[mocha]: http://visionmedia.github.com/mocha/
[expect]: https://github.com/LearnBoost/expect.js
[sinon]: http://sinonjs.org/
[double]: http://en.wikipedia.org/wiki/Test_double
