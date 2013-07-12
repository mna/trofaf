---
Date: 2012-03-27
Title: Expérimentations sur l'injection de dépendance avec node.js
Author: Martin Angers
Category: technologie
Description: L'injection de dépendance permet de découpler les différentes composantes d'une application pour en faciliter l'entretien, la testabilité, circonscrire l'impact des changements, mais aussi pour imposer une façon d'aborder la création de l'application en une aggrégation de pièces simples, à la responsabilité ciblée, et à l'API bien défini.
---

L'injection de dépendance (*dependency injection*) permet de découpler les différentes composantes d'une application pour en faciliter l'entretien, la testabilité, circonscrire l'impact des changements, mais aussi pour imposer une façon d'aborder la création de l'application en une aggrégation de pièces simples, à la responsabilité ciblée, et à l'API bien défini. Elle a comme effet secondaire, règle générale, une meilleure architecture.

Dans les langages de programmation orientés-objet statiques, tels C# ou Java, le patron (*pattern*) d'injection de dépendance prend la forme d'une interface ou d'une classe abstraite, d'une (ou plusieurs) implémentation de cette interface, et d'un *assembleur*, responsable de créer l'instance concrète d'une implémentation de l'interface, et de l'injecter dans un objet "client" de la dépendance. [Martin Fowler a écrit un brillant article sur le sujet][fowler], c'est une lecture essentielle pour approfondir ce patron.

Il existe, dans l'écosystème node, des librairies imitant le patron des langages statiques. Une simple recherche avec `npm search dependency injection` en donne un aperçu. C'est moins cette mécanique d'assemblage, relativement simple et bien rodée, qui m'intéresse dans le cas présent que les différentes façons de découpler deux modules avec Javascript, et plus spécifiquement sur la plateforme node. Dit autrement, j'expérimente sur comment obtenir les bénéfices de l'injection de dépendance, et non sur la traduction en javascript de l'implémentation traditionnelle de ce patron.

### Le laboratoire *express-boilerplate*

Avec un langage dynamique comme Javascript et son [héritage par prototype][proto], il y a d'autres façons de faire. C'est ce que j'explore dans mon projet d'expérimentation [express-boilerplate sur GitHub][eb], qui évolue continuellement et qui est né de l'intérêt à publier concrètement ma [structure de projet proposée pour une application Web avec Express][struct]. Au moment d'écrire ce billet, j'ai validé deux approches, disponibles dans deux branches git distinctes:

*	**simple-DI** : Cette approche classique, toute simple, reçoit les dépendances en paramètres, et le fichier *app.js* est responsable de l'assemblage. Par exemple, dans le module *router*:
	
		:::javascript
		module.exports = function (server, handler) {
			server.get('/', handler.renderIndex);
		};

	et dans *app.js*:

		:::javascript
		var server = require('./lib/server'),
			db = require('./lib/db'),
			handler = require('./lib/handler')(db);

		// Appeler router en lui passant ses dépendances
		require('./lib/router')(server, handler);


*	**prototype-extension** : Cette branche utilise l'extension du prototype du HTTPServer de Express pour injecter le module *config* dans le module *server*. C'est un parfait exemple de ce que j'entends quand je parle d'obtenir "les bénéfices de l'injection de dépendance", et non la retranscription du patron. Cette solution repose toutefois sur une hypothèse non négligeable: la dépendance commune de ces deux modules sur Express, et cette dépendance sur Express  n'est d'ailleurs pas découplée dans mon *boilerplate*. Mais puisqu'il s'agit d'un gabarit d'application Web visant précisément ce cadre d'application (ce *framework*), c'est une concession que je juge acceptable. Ça donne ceci:

		:::javascript
		// * * * * config.js * * * *
		express.HTTPServer.prototype.applyConfiguration = function () {
			...
		}

		// * * * * server.js * * * *
		server = express.createServer();

		// Dépendance sur le module config, injectée via cet appel
		server.applyConfiguration();

		module.exports = server;

### Versatilité vs rigueur

Cependant, ce qu'on gagne en versatilité avec les langages dynamiques, on perd en rigueur (au sens "rigueur intrinsèque au langage"). L'avantage de l'approche avec interface des langages statiques, c'est d'avoir l'assurance que la dépendance reçue expose toutes les fonctionnalités voulues - méthodes et propriétés. Avec Javascript, si on reprend mon dernier exemple, je n'ai aucune assurance que le serveur Express implémente bel et bien `applyConfiguration()`. Dans un cas simple comme celui-ci, je pourrais aisément valider avant l'appel que "applyConfiguration" est bien défini sur l'objet "server", et que c'est bel et bien une fonction. Mais dans un cas plus complexe, où plusieurs méthodes et propriétés de la dépendance sont utilisées, ça peut rapidement devenir hors de contrôle.

C'est ce qu'il manque, à mon avis, aux solutions actuelles sous Javascript. Quelque chose pour valider que le contrat attendu par le "client" soit respecté par la dépendance fournie. J'ai quelques idées sur le sujet, on verra si ça mûrira, gardez l'oeil sur [le référentiel GitHub][eb], et si vous avez commentaires et suggestions, l'espace ci-bas est là pour ça!

[fowler]: http://martinfowler.com/articles/injection.html
[proto]: http://fr.wikipedia.org/wiki/Programmation_orient%C3%A9e_prototype
[eb]: https://github.com/PuerkitoBio/express-boilerplate
[spring]: http://www.springsource.org/
[struct]: http://hypermegatop.calepin.co/structurer-une-application-web-avec-express-et-nodejs.html
