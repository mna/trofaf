---
Date: 2012-03-19
Title: Structurer une application Web avec Express et Node.js
Author: Martin Angers
Category: technologie
Description: La façon d'organiser son code, en divers répertoires et fichiers, est en bonne partie une question de préférence personnelle. Le cadre d'application Web Express n'impose aucune structure particulière, laissant le développeur totalement libre d'arranger le tout selon son inspiration du moment. Cette liberté peut parfois donner le vertige et provoquer un effet pervers: le syndrome de la page blanche.
---

** MISE À JOUR ** (22 mars 2012) : J'ai mis un exemple (simpliste, mais fonctionnel) de la structure suggérée dans ce billet sur GitHub, [express-boilerplate][exbo].

* * * *

La façon d'organiser son code, en divers répertoires et fichiers, est en bonne partie une question de préférence personnelle. Le cadre d'application Web (*Web framework*) [Express][] n'impose aucune structure particulière, laissant le développeur totalement libre d'arranger le tout selon son inspiration du moment. Cette liberté peut parfois donner le vertige et provoquer un effet pervers: le syndrome de la page blanche.

C'est probablement pourquoi [la question][stack] revient [aussi souvent][ggroups]. Bien qu'il n'y ait pas de *bonne structure* canonique, les bonnes pratiques reconnues et un peu d'expérimentation permettent d'avancer une proposition réfléchie.

### Par défaut: la structure Express

Quand je dis que Express n'impose aucune structure, c'est vrai, mais ce n'est pas toute la vérité. Il *suggère* une structure, lorsqu'on utilise l'outil de ligne de commande pour créer un cadre de départ (ex.: `express struct_express`). Ça donne ceci (avec la version 2.5.5):

    struct_express
        public
            images
            javascripts
            stylesheets
        routes
        views
        app.js

On comprend aisément que *public* contient les contenus non sécurisés, utilisés par les pages html produites par l'application et pris en charge par [Connect][] (le cadre applicatif plus générique sur lequel est construit Express) via l'intergiciel (*middleware*) de fichiers statiques.

Le répertoire *routes* porte à confusion, car à l'analyse du code qu'il contient (dans *index.js*), on n'y retrouve non pas la définition des URLs et des verbes HTTP supportés par l'application (la *route*), mais seulement l'implémentation, la logique applicative rattachée à cette route:

    :::javascript
    // Fichier /routes/index.js
    exports.index = function(req, res){
        res.render('index', { title: 'Express' })
    };

La définition de la route comme telle, elle, se retrouve à la racine de l'application, dans le fichier maître *app.js*:

    :::javascript
    // Fichier /app.js
    var express = require('express')
      , routes = require('./routes')

    // ...
    // Routes
    app.get('/', routes.index);

C'est probablement acceptable pour de petits projets ou des tests rapides, mais pour une application d'une certaine taille, une meilleure organisation du code est nécessaire. Personnellement, je veux:

*   un fichier de tête le plus bête possible, qu'il ne connaisse que les dépendances à obtenir et la façon de les assembler, sans aucune intelligence au niveau de l'implémentation.
*   des fichiers (des modules, si on adhère au vocabulaire du [CommonJS][] suivi par node) courts, simples à comprendre, respectant le [principe de responsabilité unique][srp] (*single responsibility principle*).
*   du découplage par injection de dépendance, du code facilement testable.
*   une structure respectant les bonnes pratiques, les conventions généralement acceptées par la communauté.

### Une saine organisation

Voici l'organisation que je propose, qui permet de répondre à ces exigences:

    struct_express_amelioree
        lib
            config
            db
            handler
            router
            server
        public
            css
            img
            js
        test
        views
        app.js

Quelques constats rapides:

*   Le code "serveur" se retrouve sous */lib*, et les tests automatisés sous */test*, une convention suivie par la plupart des ténors de la communauté.
*   Sous */public*, on retrouve les trois mêmes répertoires de contenu statique, mais avec des noms plus courts, simple question de préférence et quelques octets de gagnés!
*   Sous */views*, on retrouve les modèles de vues (*templates*) servant à produire les pages html. Personnellement j'utilise [jade][], mais Express supporte d'autres engins.
*   */lib* contient le coeur de l'application. Chaque sous-item a une responsabilité précise, et [grâce à la flexibilité offerte par node][nodefolders], chaque sous-item peut prendre la forme d'un fichier unique ou d'un sous-répertoire et ainsi permettre une meilleure organisation (en multiples fichiers) de ce sous-item.

Puisque c'est là l'essentiel du code serveur d'une application Web, voici en détail chacun des sous-items de */lib*:

*   **config** contient la configuration du serveur Express, soit en général la mise en place des intergiciels utilisés (*middleware*), la configuration de l'engin de vues, la configuration de la gestion des erreurs selon l'environnement d'exécution, etc. Typiquement j'utilise un seul fichier, donc il prend la forme de *config.js*.
*   **db** contient la couche d'acces aux données, de même que la définition des modèles utilisés par l'application. Certains préféreront peut-être appeler ce répertoire *models*. Personnellement j'utilise souvent [MongoDB][] et la librairie node [mongoose][], donc je structure *db* sous la forme d'un répertoire avec *index.js* pour gérer la connexion à la base de données, et un fichier distinct par modèle, chacun des modèles étant exposé via *index.js*.
*   **handler** contient la logique applicative à appliquer lors de requêtes sur les routes supportées. Je découpe habituellement les implémentations en différents fichiers, par exemple un fichier distinct pour les *handlers* des routes REST de chaque modèle. Ces implémentations sont indépendantes des routes, la définition de celles-ci étant l'affaire du...
*   **router**, qui contient la définition des routes supportées par l'application. Là aussi, je suis le même découpage que pour les *handlers*, donc un fichier contenant les routes REST d'un modèle, un fichier contenant les routes de l'interface utilisateur, etc.
*   **server** contient la création du serveur HTTP comme tel, ce qui est généralement appelé l'"app" dans les exemples d'Express, mais qui est plus spécifiquement le serveur Web (l'application étant l'ensemble des modules et leurs dépendances!). C'est souvent un module très simple, qui peut se limiter à appeler `express.CreateServer()`, donc j'utilise un seul fichier, *server.js*.

Ce qui laisse un fichier de tête *app.js* effectivement très simple et sans intelligence autre que l'assemblage des modules, l'injection des dépendances, et l'appel à `server.listen()` pour démarrer l'application Web. Ça donne une organisation saine, qui permet de respecter le patron [MVC][] où, grossièrement, le modèle est */lib/db*, la vue est */views* et le contrôleur est une combinaison de */lib/router* et */lib/handler*, le *router* jouant le rôle d'"agent messager", et le *handler* contenant la logique comme telle. D'ailleurs j'utilise un découpage semblable pour organiser le code côté client (en développement, avant de les *minifier*), dans mes fichiers sous */public/js* avec [backbone.js][backbone], mais c'est une histoire pour un autre billet.

Est-ce que votre organisation du code ressemble à ça? Utilisez-vous quelque chose de radicalement différent? Faites-en part dans les commentaires!

[express]: http://expressjs.com/
[stack]: http://stackoverflow.com/questions/9607947/how-should-i-structure-my-node-express-mongodb-app
[ggroups]: https://groups.google.com/forum/#!topic/express-js/9WrW3dxXqDs
[connect]: http://www.senchalabs.org/connect/
[srp]: http://en.wikipedia.org/wiki/Single_responsibility_principle
[commonjs]: http://www.commonjs.org/
[jade]: http://jade-lang.com/
[nodefolders]: http://nodejs.org/api/modules.html#modules_folders_as_modules
[mongodb]: http://www.mongodb.org/
[mongoose]: http://mongoosejs.com/
[mvc]: http://fr.wikipedia.org/wiki/Mod%C3%A8le-Vue-Contr%C3%B4leur
[backbone]: http://backbonejs.org/
[exbo]: https://github.com/PuerkitoBio/express-boilerplate
