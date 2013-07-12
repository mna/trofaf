---
Date: 2012-02-21
Title: Propulsé par node.js
Author: Martin Angers
Category: technologie
Description: Le développement d'une application Web est ponctué d'une multitude de choix technologiques. Pour le compte rendu Web, le coeur technologique, celui qui a orienté pratiquement tous les autres choix, n'a pourtant pas fait l'objet d'une réflexion, d'une étude comparative. Loin s'en faut. Il a plutôt été l'élément déclencheur de tout le projet.
---

Le développement d'une application Web (d'un [SaaS][1], si vous préférez) est ponctué d'une multitude de choix technologiques. 
Pour [le compte rendu Web][2], le coeur technologique, celui qui a orienté pratiquement tous les autres choix, n'a pourtant pas fait l'objet
d'une réflexion, d'une étude comparative. Loin s'en faut. Il a plutôt été l'élément déclencheur de tout le projet.

C'est au courant de sa version 0.4, quelque part en début 2011, que j'ai découvert [node.js][3] au hasard de je ne sais plus trop quel heureux furetage.
Après quelques expérimentations et le visionnement de [la mythique présentation][4] de son créateur, Ryan Dahl, je savais que je devais construire quelque
chose avec cette merveille. Dire que le compte rendu Web est un prétexte pour utiliser node.js serait exagéré, mais j'avais effectivement le choix du
moteur technologique avant l'idée d'application.

### Entrées/sorties événementielles, rien ne bloque!

C'est amplement documenté, donc je passerai rapidement, mais la *révolution node* est que tous les appels qui, typiquement, sont bloquants (interrompent
le traitement en attendant que l'appel soit complété - comme à la lecture d'un fichier sur disque, par exemple) sont gérés de façon asynchrone, permettant ainsi au
traitement de se poursuivre et d'être notifié par une procédure de rappel (un *callback*) lorsque l'action demandée est complétée. Cette approche permet d'écrire du code 
qui s'exécute sur un seul fil d'exécution (un *thread*), simplifiant drastiquement la gestion de concurrence d'accès, et offrant une performance et une capacité 
de croissance (*scalability*) impressionnantes.

Ce n'est [ni le seul][5], [ni le premier][6] projet adoptant cette philosophie, mais c'est possiblement celui qui le fait le mieux et qui rend l'utilisation d'opérations
synchrones, bloquantes, l'exception à spécifier explicitement plutôt que la règle. Cet exemple simpliste démontre la facilité avec laquelle on peut coder un serveur HTTP qui
lit et retourne le contenu d'un fichier au client, le tout sans bloquer le traitement:

	:::javascript
	var http = require('http'),
    	fs = require('fs');

	http.createServer(function (req, res) {
		res.writeHead(200, {'Content-Type': 'text/plain'});
		fs.readFile('helloworld', function (err, data) {
			if (err) throw err;
			res.write(data);
			res.end();
		});
	}).listen(1337, "127.0.0.1");

Cet exemple démontre une autre particularité de node: contrairement à la plupart des autres technologies où le serveur HTTP est un logiciel externe à l'application 
(pensons à Apache ou IIS, par exemple), ici le serveur est un membre à part entière, avec tout le contrôle que cela implique.

L'exemple démontre aussi l'évidence pour l'oeil averti: le langage utilisé pour coder avec node est javascript (d'où le `.js` dans node.js). Ce n'est pas un choix innocent, javascript étant 
naturellement adapté pour une approche par procédure de rappel, et n'ayant aucune librairie existante (à l'origine) pour faire de l'entrée/sortie, il n'y avait
donc pas ce danger d'utiliser une librairie bloquante. Les fonctionnalités d'entrée/sortie, les librairies de gestion de protocoles internet et autres ont pu
être développées à partir d'une feuille vierge, avec dès le départ cette notion d'asynchronisme, sans avoir à travestir de l'existant. Et l'engin d'exécution
du javascript est [le fameux V8 utilisé dans le fureteur Chrome de Google][7] - autre morceau essentiel du puzzle, car il n'y a pas si longtemps, javascript n'avait
aucun environnement d'exécution suffisamment performant pour oser rêver une utilisation *sérieuse* côté serveur.

### L'esprit de communauté

Tout ça est bien beau, mais la vraie force de node réside dans sa bouillonnante communauté concentrée sur le portail de partage de code source libre [GitHub][8].
Se sont greffés, dans le sillon de Dahl, de fabuleux talents tels:

*	[T.J. Holowaychuk][11], auteur principal de la couche d'intergiciel (*middleware layer*) 
	pour serveur HTTP [Connect][9], du populaire cadre d'application Web (*framework*) [Express][12], et de l'outil de test [BDD][10] [Mocha][13], entre autre.
*	[Tim Caswell][14], auteur du gestionnaire de versions de node [nvm][15], qui permet d'installer différentes versions sur un poste et d'activer celle que l'on veut pour
	une session de ligne de commande bash (Linux ou Mac seulement), et créateur du blogue [How To Node][16].
*	[Mikeal Rogers][17], auteur de la simplissime librairie de requêtes HTTP [Request][18], et organisateur de l'événement [NodeConf][19].
*	Et pour relier ces îlots de génie en un tout symbiotique, il y a [Isaac Schlueter][20], auteur du gestionnaire de librairie [npm][21], maintenant déployé par défaut
	avec node. L'installation d'une librairie pour node est aussi triviale que d'écrire sur la ligne de commande `npm install <nom_librairie>`. À ce jour,
	[plus de 7000 librairies][22] sont disponibles via npm. Isaac a récemment [remplacé Ryan Dahl comme maître du code node][24].

Ce ne sont que quelques-uns des contributeurs exceptionnels de cet écosystème. Même [Microsoft n'a pu résister à la vague node][23], ayant participé au 
port de node sur la plateforme Windows.

C'est donc sur cette base solide et hyper motivante que repose le [Compte Rendu Web][2]. Propulsé par node.js, à plus d'un niveau!

[1]: http://fr.wikipedia.org/wiki/Logiciel_en_tant_que_service "Logiciel en tant que service"
[2]: http://www.compterenduweb.com/
[3]: http://nodejs.org/
[4]: http://jsconf.eu/2009/video_nodejs_by_ryan_dahl.html
[5]: http://twistedmatrix.com/trac/ "Python Twisted"
[6]: http://rubyeventmachine.com/ "Ruby Event Machine"
[7]: http://code.google.com/p/v8/ "Google V8"
[8]: https://github.com/
[9]: http://www.senchalabs.org/connect/
[10]: http://fr.wikipedia.org/wiki/Behavior_Driven_Development "Behavior-Driven Development"
[11]: http://tjholowaychuk.com/
[12]: http://expressjs.com/
[13]: http://visionmedia.github.com/mocha/
[14]: https://github.com/creationix
[15]: https://github.com/creationix/nvm
[16]: http://howtonode.org/
[17]: http://www.mikealrogers.com/
[18]: https://github.com/mikeal/request
[19]: http://www.nodeconf.com/
[20]: https://github.com/isaacs
[21]: https://github.com/isaacs/npm
[22]: http://search.npmjs.org/
[23]: http://blog.nodejs.org/2011/06/23/porting-node-to-windows-with-microsoft%E2%80%99s-help/
[24]: http://venturebeat.com/2012/01/30/dahl-out-mike-drop/
