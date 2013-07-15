---
Date: 2012-03-01
Title: npm shrinkwrap: Comment contrôler ses dépendances
Author: Martin Angers
Category: technologie
Description: Utiliser une librairie existante permet d'ajouter rapidement des fonctionnalités à une application, et de concentrer ses efforts sur les nouveautés, les spécificités de son projet, plutôt qu'à réinventer la roue. Cependant, chaque librairie utilisée devient une dépendance, et une dépendance est un risque. Un risque car notre code dépend maintenant de ce corps étranger sur lequel on a généralement peu ou pas de contrôle. Les meilleures pratiques pour limiter ce risque ont évolué rapidement sur la plateforme node.js ces derniers mois.
---

Utiliser une librairie existante permet d'ajouter rapidement des fonctionnalités à une application, et de concentrer ses efforts sur les nouveautés, les spécificités de son projet, plutôt qu'à réinventer la roue. Cependant, chaque librairie utilisée devient une dépendance, et une dépendance est un risque. Un risque car notre code dépend maintenant de ce corps étranger sur lequel on a généralement peu ou pas de contrôle. Les meilleures pratiques pour limiter ce risque ont évolué rapidement sur la plateforme node.js ces derniers mois.

> "A dependency upon a package is a dependency upon everything within the package. When a package changes, and its release number is bumped, all clients of that package must verify that they work with the new package - even if nothing they used within the package actually changed."  
> - Robert C. Martin

### Le cas général: package.json

Le premier niveau de contrôle est fourni par npm via le fichier de méta-données [package.json][json]. Celui-ci offre quatre clefs de configuration pour gérer les dépendances de l'application:

1. **dependencies**: un dictionnaire où la clef est le nom de la librairie, et la valeur est une version, qui peut s'exprimer de plusieurs façons, j'y reviens dans un instant.
2. **devDependencies**: un autre dictionnaire qui prend la même forme que pour *dependencies*, mais qui contient les dépendances requises seulement en mode développement, généralement les librairies de test automatisé ou de génération de documentation. Ces dépendances ne sont pas installées par défaut, mais elles peuvent l'être avec la commande `npm install <librairie> -dev`.
3. **bundledDependencies**: un tableau (*array*) complémentaire à *dependencies* qui liste les librairies qui devront être publiées **avec** la librairie décrite par ce *package.json*. C'est une configuration utilisée en de rares cas, dont npm lui-même, puisqu'il ne peut dépendre sur npm pour déployer ses dépendances! Règle générale, il est recommandé de ne *pas* utiliser cette approche - qui force à redéployer ce regroupement de librairies à chaque changement dans une dépendance - et plutôt de laisser npm résoudre le tout dynamiquement.
4. **engine**: un dictionnaire un peu particulier, en ce sens qu'il est l'équivalent de *dependencies*, mais pour spécifier les versions de node ciblées (et possiblement les versions de npm).

Chaque dictionnaire permet de spécifier une ou des versions permises pour chaque dépendance. L'information de version peut être spécifiée de nombreuses façons:

*	Version spécifique (ex.: "1.2.4" ou "=1.2.4")
*	Version minimum ou maximum (ex.: "<1.2.4", ">=1.2.4")
*	Toute version ("*" ou "")
*	Une étendue (sous la forme ">=1.2.4 <1.3.0", ou "1.2.4 - 1.3.2", qui est une étendue *inclusive*, donc équivalent à ">=1.2.4 <=1.3.2", ou encore "1.2.x" où le "x" peut être n'importe quel chiffre, donc dans cet exemple, correspond à ">=1.2.0 <1.3.0", et "1.x" correspond à ">=1.0.0 <2.0.0")
*	Une étendue spécifiée par un tilde ("~1.2.4" qui correspond à ">=1.2.4 <1.3.0", alors que "~1.2" correspond à ">=1.2.0 <2.0.0")

Il est de plus possible de combiner différentes étendues grâce à l'opérateur `||` qui est un "ou" logique (le "et" est l'opérateur par défaut, soit l'espace entre deux versions). Ainsi, cet exemple est valide: `"<1.1.0 || >=3.1.0 <3.3.0"`. npm supporte même, en lieu et place d'un critère de version, une adresse URL pointant vers une archive de type [tar][] (un *tarball*) contenant la librairie voulue et un fichier package.json, ou un référentiel (*repository*) [git][].

Les versions de librairies doivent se conformer au [versionnage sémantique][semver], [légèrement adapté][nodesemver] dans le contexte de l'écosystème node. L'essentiel du versionnage sémantique est que toute version comporte trois chiffres: M.m.r (**M**ajeure, **m**ineure, **r**ustine ou *patch*). Il est possible d'ajouter un numéro de génération (*build*) sous la forme M.m.r-gen, ou un libellé quelconque (*tag*), par exemple 1.2.0beta. La version sans libellé est considérée plus grande que celle avec libellé ("1.2.0" > "1.2.0beta"). Les principales règles à suivre, et elles sont essentielles à la bonne cohabitation des dépendances, sont les suivantes:

*	Toute nouvelle publication doit avoir une nouvelle version, supérieure à la précédente (c'est l'évidence).
*	Le numéro de rustine doit être incrémenté si les seuls changements sont des corrections dans l'implémentation interne, sans impact sur l'[API][] exposé.
*	Le numéro mineur doit être incrémenté si de nouvelles fonctionnalités sont introduites sans bris de compatibilité dans l'API public, ou si une fonctionnalité de l'API est en délestage, ou si des changements substantiels sont introduits dans l'implémentation interne. Le numéro de rustine doit être réinitialisé à 0 lors d'une incrémentation de numéro mineur.
*	Le numéro majeur doit être incrémenté si un bris de compatibilité est introduit dans l'API public. Le numéro mineur et celui de rustine doivent être réinitialisés à zéro.

C'est pourquoi il n'est pas rare de voir des dépendances définies par une étendue sur la rustine ou même le numéro mineur ("1.2.x" ou "1.x", par exemple), puisque ça garantit, si le versionnage sémantique est honoré par l'auteur de la librairie, que les versions dans ces étendues respecteront toutes au minimum l'API utilisé.

### Problématique et piste de solution

Le problème avec cette approche, c'est que bien que les versions des dépendances directes de l'application soient contrôlées, les dépendances de ces dépendances (et ainsi de suite dans la hiérarchie) sont hors de notre contrôle. En fait, elles sont contrôlées par le package.json de leur librairie "parent". Ainsi, il est possible qu'en développement, on ait cette hiérarchie:

    appA (v0.1.0)
     - libB (v0.0.1)
       - libC (v0.1.0)

Et que lors du déploiement en production, via un `npm install` tout frais sur notre *appA* rigoureusement testée:

    appA (v0.1.0)
     - libB (v0.0.1)
       - libC (v0.2.0)

Sans crier gare, voilà que la version 0.2.0 de *libC* fait son entrée en scène, en production, sans jamais avoir été testée dans le contexte de notre application.

C'est ce type de problème qui a mené Mikeal Rogers à écrire [ce billet][nmingit] sur une nouvelle façon de faire: inscrire le répertoire "node_modules" (qui contient toutes les dépendances) dans le référentiel de code source git. L'idée est intéressante, l'application est développée et testée avec un jeu de dépendances, le tout est inscrit dans git, donc un simple `git clone` du référentiel source installe l'application et toutes ses dépendances, figées récursivement lors du dernier `git commit`.

Attention, on parle ici d'*applications à déployer*, **pas** de *librairies réutilisables*. Ces dernières ne devraient pas figer leurs dépendances de façon aussi drastique. Mikeal aborde le sujet dans son article sous l'angle de la responsabilité des tests d'intégration, j'ajouterais un autre argument pourquoi les librairies réutilisables devraient s'en tenir aux étendues de version du package.json et ne pas figer dans le temps leurs dépendances: elles ne savent pas quand elles seront déployées en production. Pour une application, on fige les versions des dépendances à un moment précis dans le temps, généralement vers la fin du développement, avant la gamme d'essais précédant la mise en production. Si les librairies réutilisables figent leurs dépendances au moment de leur *publication*, elles privent potentiellement leurs applications "clientes" de correctifs importants dans leurs propres dépendances.

Pour reprendre l'exemple précédent de l'*appA*, il est normal et même souhaitable que *libB* permette à sa dépendance *libC* d'évoluer dans les limites de la compatibilité de son API. Cependant, l'exemple souligne l'importance d'avoir une mécanique permettant aux applications "déployables" de garder la maîtrise du processus.

### L'emballage cadeau: shrinkwrap

Bien que théoriquement solide, l'approche des dépendances dans git a quelques défauts, principalement au niveau des [dépendances binaires][bindeps] (voir dans l'article lié la section "*Why not just check node_modules into git?*"). Il est difficile de bien déterminer ce qui doit être exclu de git (seul le code source servant à reconstituer le binaire devrait être inclus), et tel que documenté dans l'article référencé, il faut alors regénérer les binaires dans l'environnement cible, lors d'un déploiement, ce qui n'est pas sans risque. L'option d'inscrire même les binaires dans git a été envisagée, mais il y a alors le risque d'erreur humaine (si le code est changé sans recompiler le binaire).

La solution proposée est cette nouvelle commande de npm, [lancée sur le blogue de node][bindeps] (et qui a eu étonnamment peu d'écho jusqu'ici, si je me fie à mes recherches google), *shrinkwrap*! Son usage est très simple, l'exécution de `npm shrinkwrap` dans le répertoire d'une application avec un package.json conforme (les dépendances doivent y être correctement listées, et aucune librairie superflue ne doit être installée dans "node_modules") génère un fichier "npm-shrinkwrap.json" qui contient la version *actuellement installée* de chacune des dépendances, de façon récursive. Lorsque ce fichier est présent, npm l'utilise plutôt que les package.json pour contrôler l'installation des librairies. Citation de [Dave Pacheco][pacheco] de Joyent - la société derrière node.js:

> When "npm install" installs a package with a npm-shrinkwrap.json file in the package root, the shrinkwrap file (rather than package.json files) completely drives the installation of that package and all of its dependencies (recursively).

À titre d'exemple, voici un extrait de mon fichier shrinkwrap pour le [Compte Rendu Web][crw]:

    {
	  "name": "crw",
	  "version": "0.1.0",
	  "dependencies": {
	    "debug": {
	      "version": "0.5.0"
	    },
	    "express": {
	      "version": "2.5.8",
	      "dependencies": {
	        "connect": {
	          "version": "1.8.5",
	          "dependencies": {
	            "formidable": {
	              "version": "1.0.9"
	            }
	          }
	        },

Il est encore possible de demander à npm quelles librairies ne sont pas à jour avec `npm outdated`, ce qui donne dans mon cas:

    mongoose@2.5.10 ./node_modules/mongoose current=2.5.9
    mime@1.2.5 ./node_modules/express/node_modules/connect/node_modules/mime current=1.2.4

Si je veux mettre mes dépendances à jour, c'est toujours possible avec `npm update`, mais tant que je ne recrée pas le fichier *shrinkwrap*, si je réexécute `npm install`, je reviens aux vieilles versions spécifiées dans ce fichier. Pour conserver les versions à jour, je dois refaire `npm shrinkwrap` suite au `npm update`. Autrement dit, du moment que *shrinkwrap* est utilisé, les versions installées par npm sont entièrement sous mon contrôle.

### Brèche de sécurité dans le registre

Un mot en terminant sur le [registre][], je mentionnais dans mon [dernier article sur npm][moi] qu'il était nécessaire de se créer un compte pour utiliser certaines commandes. Et bien une brèche de sécurité a été découverte (et corrigée). Un courriel a été envoyé par Isaac à tous ceux qui avaient fourni une adresse valide, mais je sais qu'il a eu des problèmes avec GMail et certains n'ont pas reçu l'information, donc si c'est votre cas, assurez-vous de consulter ce [gist][] qui reprend le contenu du courriel et les étapes à suivre pour sécuriser votre compte.

[json]: http://npmjs.org/doc/json.html
[tar]: http://fr.wikipedia.org/wiki/Tar_(informatique)
[git]: http://fr.wikipedia.org/wiki/Git
[semver]: http://semver.org/
[nodesemver]: http://npmjs.org/doc/semver.html
[api]: http://fr.wikipedia.org/wiki/Interface_de_programmation
[nmingit]: http://www.mikealrogers.com/posts/nodemodules-in-git.html
[bindeps]: http://blog.nodejs.org/2012/02/27/managing-node-js-dependencies-with-shrinkwrap/
[crw]: http://www.compterenduweb.com/
[moi]: http://hypermegatop.calepin.co/npm-la-base-essentielle-pour-debuter-avec-nodejs.html
[gist]: https://gist.github.com/2001456
[pacheco]: http://blog.nodejs.org/author/davepacheco/
[registre]: http://search.npmjs.org/
