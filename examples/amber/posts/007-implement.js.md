---
Date: 2012-04-02
Title: implement.js: typage fort et Javascript
Author: Martin Angers
Category: technologie
Description: L'injection de dépendance avec Javascript a comme conséquence de ne pouvoir assumer que les fonctionnalités offertes par l'instance reçue seront celles attendues.
---

L'injection de dépendance (*dependency injection*) avec Javascript a comme conséquence de ne pouvoir assumer que les fonctionnalités offertes par l'instance reçue seront celles attendues. Dans les langages statiques, l'injection est généralement basée sur une interface, qui assure au module "client" la disponibilité d'un ensemble de fonctionnalités. Le contrat imposé par l'interface est assurément respecté. Dans un langage dynamique comme Javascript, il n'y a rien de tel.

Ce qui laisse deux options aux développeurs Javascript: l'acte de foi ou la validation manuelle de l'objet reçu. J'en parlais dans [mon dernier billet][deps], j'avais l'impression qu'il manquait quelque chose à ce niveau.

### implement.js

C'est là qu'entre en scène ma nouvelle librairie offerte en logiciel libre, [implement.js][impl] (comme vous pouvez le constater, *je blog in French, mais je code in English*). Elle permet de valider qu'un objet respecte une interface prédéfinie. Ou que les paramètres passés à une fonction sont bel et bien du type attendu. Ce sont là les deux principales fonctionnalités exposées par la librairie.

Je ne reprendrai pas ici [le contenu du *readme* sur GitHub][readme], je vous suggère de le lire et de m'avertir de toute information manquante ou ambigüe. Il s'agit de ma première contribution "significative" au merveilleux monde du logiciel libre, et je le souhaite, pas la dernière! J'ai aussi ajouté une branche dans mon projet expérimental [express-boilerplate][eb] pour intégrer l'utilisation de *implement.js*, qui peut servir d'exemple d'utilisation, tout comme [le répertoire *examples*][examples] de *implement.js*.

### Essayez et participez!

Les prochaines étapes sont assez simples, je vous invite à essayer la librairie, si vous croyez qu'elle peut être utile dans vos projets. N'hésitez pas à la faire partager avec vos collègues développeurs, et à [ouvrir un rapport d'anomalie][issue] au besoin. Encore mieux, si vous le pouvez, soumettez un correctif (un *pull request*). Je prévois ajouter la documentation de l'API avec [JSDoc][] sous peu, et possiblement rendre la librairie disponible dans le fureteur.

[deps]: http://hypermegatop.calepin.co/experimentations-sur-linjection-de-dependance-avec-nodejs.html
[impl]: https://github.com/PuerkitoBio/implement.js
[issue]: https://github.com/PuerkitoBio/implement.js/issues
[jsdoc]: http://en.wikipedia.org/wiki/JSDoc
[readme]: https://github.com/PuerkitoBio/implement.js#readme
[eb]: https://github.com/PuerkitoBio/express-boilerplate
[examples]: https://github.com/PuerkitoBio/implement.js/tree/master/examples
