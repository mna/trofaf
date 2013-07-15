---
Date: 2012-03-01
Title: npm: la base essentielle pour débuter avec node.js
Author: Martin Angers
Category: technologie
Lang: fr
Description: La plateforme node.js est volontairement limitée, en son coeur, aux fonctionnalités les plus fondamentales. Elle mise donc sur les contributions de sa communauté pour enrichir le noyau et repousser les limites du possible. Avec près de 8000 librairies à ce jour, il y a sérieux danger de chaos. Et si on y retrouve plutôt quelque chose comme une belle organisation, c'est en bonne partie grâce au gestionnaire de librairies npm. 
---

La plateforme node.js est [volontairement limitée][userland], en son coeur, aux fonctionnalités les plus fondamentales. Elle mise donc sur les contributions - souvent époustouflantes, [j'en parle ici][billetnode] - de sa communauté pour enrichir le noyau et repousser les limites du possible.

Avec près de 8000 librairies pullulant dans son écosystème, il y a sérieux danger de chaos. Et si on y retrouve plutôt quelque chose comme une belle organisation, c'est en bonne partie grâce à la colonne vertébrale qui unit et soutient cette terre fertile, le gestionnaire de librairies [npm][]. Présentation, trucs et astuces dans les lignes qui suivent.

### La base: installation, désinstallation, indicateur global

*npm* lui-même s'installe de façon fort simple depuis quelques versions: il est livré avec node. Il s'agit donc d'installer node d'une des [nombreuses][installnode] [façons][nvm] [possibles][n].

Ensuite, toute librairie publiée dans [le registre de npm][registry] est à portée de main. Pour installer une librairie, par exemple le gestionnaire de versions *n*:

    :::Bash shell scripts
    martin@LilDevil:~/sublime-wrkspc/crw$ npm install n
	npm http GET https://registry.npmjs.org/n
	npm http 200 https://registry.npmjs.org/n
	npm http GET https://registry.npmjs.org/n/-/n-0.7.0.tgz
	npm http 200 https://registry.npmjs.org/n/-/n-0.7.0.tgz
	n@0.7.0 ./node_modules/n

Seul bémol à propos de cette démonstration rayonnante de simplicité: la commande ne respecte pas [la règle de silence de UNIX][unixrules], qui veut que si tu n'as rien de spécial à dire, tu ne dis rien. En effet, sur UNIX et ses déclinaisons, on s'attend à ce qu'une commande qui fait ce qu'on lui demande sans problème n'affiche aucun message. C'est mineur, mais c'est une des philosophies importantes à respecter pour bien s'intégrer à l'ensemble.

Une fois installée, la librairie est au bon endroit pour être retrouvée par node et utilisable dans notre code, via la commande habituelle `var librairie = require('librairie');`. Si on veut installer non seulement la librairie et ses dépendances pour qu'elle s'exécute correctement (l'installation par défaut), mais aussi ses dépendances de développement, par exemple pour exécuter les tests automatisés de la librairie, alors on doit spécifier l'indicateur -dev (`npm install n -dev`). Pour désinstaller, on s'en doute, c'est `npm uninstall n` (ou un alias équivalent, soit `npm remove` ou `npm rm`).

Autres commandes de base toujours utiles, pour connaître la version actuelle de npm:

    :::Bash shell scripts
    martin@LilDevil:~/sublime-wrkspc/crw$ npm -v
    1.1.4

Et pour lister les librairies installées dans le répertoire courant (dans cet exemple, les librairies utilisées par le [Compte Rendu Web][crw] à ce point-ci de son développement):

    :::Bash shell scripts
    martin@LilDevil:~/sublime-wrkspc/crw$ npm ls
	crw@0.1.0 /home/martin/sublime-wrkspc/crw
	├── debug@0.5.0 
	├─┬ express@2.5.8 
	│ ├─┬ connect@1.8.5 
	│ │ ├── formidable@1.0.9 
	│ │ └── mime@1.2.5 
	│ ├── mime@1.2.4 
	│ ├── mkdirp@0.3.0 
	│ └── qs@0.4.2 
	├─┬ jade@0.20.3 
	│ ├── commander@0.5.2 
	│ └── mkdirp@0.3.0 
	├── less@1.2.2 
	├─┬ mocha@0.14.0 
	│ ├── commander@0.5.2 
	│ ├── diff@1.0.2 
	│ └── growl@1.5.0 
	├─┬ mongoose@2.5.9 
	│ ├── hooks@0.1.9 
	│ └── mongodb@0.9.7-3-5 
	├── request@2.9.153 
	├── should@0.6.0 
	└── uglify-js@1.2.5 

`npm list` est équivalent, et `npm ll` (ou `npm la`) offre des informations plus détaillées. Par défaut, npm conserve toutes les dépendances localement, dans un sous-répertoire "node_modules" du répertoire courant. Ainsi, comme chaque librairie installée vient avec ses propres dépendances dans son propre sous-répertoire "node_modules", une application node peut utiliser les librairies A et B, chacune ayant une dépendance sur C, mais A et B peuvent utiliser une version différente de C sans problème. On en voit un exemple concret sous la librairie *express*, qui dépend directement sur *mime*, alors que sa dépendance *connect* dépend aussi sur *mime*, mais dans une version distincte, et ce sans conflit:

    :::Bash shell scripts
    ├─┬ express@2.5.8 
	│ ├─┬ connect@1.8.5 
	│ │ ├── formidable@1.0.9 
	│ │ └── mime@1.2.5 
	│ ├── mime@1.2.4 
	│ ├── mkdirp@0.3.0 
	│ └── qs@0.4.2 

C'est une qualité assez exceptionnelle parmi les environnements de développement, que [Mikeal Rogers n'a pas manqué de souligner][mikeal]:

> Node’s local module support accomplishes what no other platform I know of has done, it allows for two dependencies to require entirely different versions of the same dependency without caveats and unforeseen failures.

Parfois, on peut préférer une installation globale d'une librairie, par exemple lorsqu'il s'agit d'une application console (*command-line interface (CLI)*) que l'on veut rendre disponible partout. npm supporte ce mode via l'indicateur `-g` ou l'équivalent `--global`. Ainsi, pour installer globalement: `npm install -g <librairie>`, pour lister les librairies installées globalement: `npm ls -g`, et ainsi de suite, les mêmes commandes s'appliquent mais avec une portée globale. À noter que ce n'est **pas** l'option recommandée pour une librairie que l'on veut utiliser dans le code d'une application node. L'approche locale à cette application est à privilégier dans ce cas (on verra pourquoi dans un futur article sur la gestion des dépendances).

### Le registre au bout des doigts

J'ai mentionné [le registre de npm][registry], qui permet de rechercher avec des mots-clefs dans les nombreuses librairies disponibles. Or, il n'est pas nécessaire d'aller sur le site Web pour effectuer ces recherches. Avec npm, on a ce registre au bout des doigts, exploitable avec de simples commandes.

    :::Bash shell scripts
    martin@LilDevil:~$ npm search /^jscov
	NAME        DESCRIPTION                    AUTHOR    DATE              KEYW
	jscoverage  jscoverage module for node.js  =kate.sf  2012-02-28 02:11

Ainsi, `npm search` suivi de mots à rechercher est équivalent à la recherche par le site Web. Il est également possible, tel que démontré dans l'exemple, d'utiliser une expression régulière en commençant le mot à rechercher par `/`. Pour afficher les informations détaillées d'une librairie spécifique - par exemple ses dépendances et ses auteurs, c'est la commande `npm view <librairie>`.

Pour profiter des fonctionnalités non anonymes, il faut bien sûr se créer un compte utilisateur. C'est tout simple, avec `npm adduser`, on peut choisir un nom d'utilisateur, un mot de passe, et associer le tout à une adresse courriel. Ensuite il est possible d'utiliser une des commandes méconnues qui a pourtant un potentiel intéressant pour découvrir les librairies de qualité: `npm star <librairie>`. C'est l'équivalent d'un *j'aime* de Facebook, et c'est une information visible dans le registre. Elle gagnerait à être davantage utilisée et mieux exploitée, car plus l'écosystème grossit, plus il devient difficile de discerner les librairies de calibre production, maintenues et codées de façon professionnelle, des projets de fin de semaine avec plus ou moins d'avenir. C'est une problématique [qui commence à faire beaucoup jaser][trouble] dans la communauté, et une solution - ou du moins une amélioration - serait sur les planches à dessin du côté de l'équipe de npm.

Dans un prochain article j'aborderai les commandes davantage liées à une librairie spécifique, son fichier `package.json`, et les bonnes pratiques (toujours en évolution!) pour la gestion des dépendances d'une application node.

[userland]: https://github.com/joyent/node/wiki/node-core-vs-userland
[billetnode]: http://hypermegatop.calepin.co/propulse-par-nodejs.html
[npm]: http://npmjs.org/
[installnode]: https://github.com/joyent/node/wiki/Installation
[nvm]: https://github.com/creationix/nvm
[n]: https://github.com/visionmedia/n
[registry]: http://search.npmjs.org/
[crw]: http://www.compterenduweb.com/
[mikeal]: http://www.mikealrogers.com/posts/nodemodules-in-git.html
[unixrules]: http://www.faqs.org/docs/artu/ch01s06.html
[privateregistry]: http://npmjs.org/doc/registry.html
[trouble]: http://mikkel.hoegh.org/blog/2011/12/20/trouble-in-node-dot-js-paradise-the-mess-that-is-npm/
