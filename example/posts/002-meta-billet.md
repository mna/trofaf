---
Date: 2012-02-29
Title: Méta-billet: un mot sur le Calepin!
Author: Martin Angers
Category: technologie
Description: Il existe de nombreux moteurs de blogue gratuits sur internet. Pourquoi avoir jeté l'ancre sur ce discret et modeste Calepin?
---

Il existe de nombreux moteurs de blogue gratuits sur internet, tels [Wordpress][wp], [Blogger][] et [Tumblr][]. Pourquoi avoir jeté l'ancre sur ce discret et modeste [Calepin][]? Peut-être un peu pour ça, tiens, pour me tenir loin des tumultes de ces méga-centres urbains du Web. Un peu aussi pour ce *design* minimaliste, centré sur le contenu, sur les mots. Et un peu pour deux qualités franchement moins romantiques: *dropbox* et *markdown*.

### Mes billets, mon coffret

La plupart des [blogiciels][blogiciel] utilisent une base de données quelconque pour stocker leurs articles. Wordpress, par exemple, repose sur une base de données [MySql][wpmysql]. Pour éviter à l'utilisateur l'horrible expérience d'inscrire son texte directement dans la base de données via des commandes SQL, ces blogiciels offrent, généralement sous la forme d'une page Web privée, des outils d'édition maison, permettant d'insérer, modifier ou supprimer les articles de façon relativement conviviale. C'est bien, c'est même un aspect important de ce qu'il est convenu d'appeler le [Web 2.0][web20], mais c'est une approche qui comporte quelques irritants:

*	Le contenu doit être saisi dans l'outil d'édition fourni. Ou importé d'un format tiers, si l'outil d'édition le permet. Bon nombre de formats sont théoriquement "compatibles" avec Wordpress, moyennant l'installation de tel ou tel [plugiciel][] plus ou moins supporté, maintenu et fonctionnel, mais c'est quand même une étape de plus, donc de trop.
*	Les billets sont conservés dans l'inutilisable format de la base de données. Bien sûr, on peut se tourner vers une fonctionnalité d'exportation du contenu, question d'obtenir, par exemple, un sympathique fichier XML (oui, c'est du sarcasme), le tout avec plus ou moins de contrôle sur ce qui sera exporté.

On parle pourtant de blogue ici, d'articles écrits pour être lus. Je peux imaginer assez aisément bon nombre d'auteurs écrire d'abord leurs billets dans le traitement de texte de leur choix, pour ensuite le transcrire, le copier ou l'importer tant bien que mal dans le blogiciel. Le Calepin prend une approche plus respectueuse du texte: mes billets restent dans leur fichier d'origine, bien au chaud et en sécurité dans *mon* coffret dans les nuages, *mon* dropbox. Le Calepin n'a besoin que d'un accès autorisé à un répertoire spécifique de dropbox, et à mon signal, il lit le contenu de ce répertoire et construit le blogue à partir de ce qu'il y trouve.

Je ferme mon blogue? Je largue les amarres? Mes billets sont toujours confortablement installés dans mon nuage. Et en tout temps, ils demeurent totalement lisibles et modifiables par l'outil de *mon* choix. Ce qui m'amène à parler du format de ces articles.

### Minimaliste jusque dans la syntaxe

Le format, la langue que comprend le Calepin, est le [markdown][]. C'est pratiquement un anti-format, en ce sens que le *markdown* est d'abord et avant tout du texte brut, sans flafla, sans encodage binaire ou enrobage XML, pouvant être lu et modifié avec le plus humble des éditeurs de texte. Là où le HTML - un langage de type *markup*, contraste intéressant - ajoute des balises dans le texte, rendant la lecture de la source difficile ou à tout le moins désagréable sans traitement préalable par un fureteur, le *markdown* utilise un nombre limité de légers artifices pour permettre d'identifier, par exemple, des en-têtes, une emphase, ou des points de forme. Ces artifices se traduisent aisément en code HTML ou une multitude d'autres formats via différents outils (dont le Calepin) mais même sous leur forme originelle, l'effet recherché est tout à fait perceptible (on peut d'ailleurs visualiser la source de chaque article sur Calepin en remplaçant le ".html" dans l'adresse par ".txt"). Un exemple:

	:::markdown
	## Ceci est un en-tête ##

	Et voici une *emphase*, et _une autre_ avec une syntaxe alternative.

	*	Ceci est un point de forme.
	*	Et un autre.

	> Et voici une citation, qui rappelle la syntaxe utilisée dans les 
	> courriels lorsqu'on répond en citant le message original.

D'une simplicité géniale. Plus de problème de compatibilité entre Mac, Windows, Linux, Android ou iOS, c'est du texte encodé en UTF-8, tout ce qu'il y a de plus standard aujourd'hui. En fait, certains en ont même fait leur [format portable à utiliser pour toute documentation][mdhn], au détriment des Word et LibreOffice de ce monde. Je considère fortement permettre cette syntaxe pour le [Compte Rendu Web][crw], puisqu'elle combine expressivité, simplicité et rapidité, trois caractéristiques essentielles pour la prise de note en temps réel.

### De la personnalité, une signature

Le Calepin enrobe le coffret dans les nuages de *dropbox* et la syntaxe pure et naturelle du *markdown* d'une fine couche de gestion à l'attention de l'auteur du blogue. À l'image de tout le reste, la page d'administration est dénuée de toute exhubérance, en trois courtes étapes tout est fait:

![Administration du Calepin](http://dl.dropbox.com/u/21605004/CalepinAdmin.jpg)

1.	L'authentification à son compte *dropbox*, qui sert par le fait même d'authentification au blogue.
2.	La configuration du blogue, soit le titre, l'option d'y ajouter le gestionnaire de commentaires qui est [Disqus][] et c'est tout, l'option du lien social qui est [Twitter][] et c'est tout, et le nom de domaine si on ne veut pas celui du Calepin par défaut.
3.	Le bouton de publication.

Voilà. Il y a quelques options mineures au niveau des méta-données de chaque article et une configuration pour le site, principalement des valeurs par défaut, mais rien de renversant, et ce *par choix*.  Pas de blogoliste (*blog roll*), de nuage de mots-clefs (*tag cloud*), de catalogue de thèmes dans tous les tons de marron, de barre de partage pour trente-douze gazillions de sites sociaux. C'est la sympathique personnalité du Calepin, minimaliste jusqu'au bout!

[blogiciel]: http://www.oqlf.gouv.qc.ca/ressources/bibliotheque/dictionnaires/terminologie_blogue/blogiciel.html
[wpmysql]: http://codex.wordpress.org/FAQ_Developer_Documentation#Why_does_WordPress_only_support_MySQL.3F_What_about_DB_abstraction.3F
[web20]: http://fr.wikipedia.org/wiki/Web_2.0
[wp]: http://wordpress.org/
[blogger]: http://blogger.com/
[tumblr]: https://www.tumblr.com/
[plugiciel]: http://www.oqlf.gouv.qc.ca/ressources/bibliotheque/dictionnaires/internet/fiches/1299146.html
[calepin]: http://calepin.co/
[markdown]: http://daringfireball.net/projects/markdown/basics
[mdhn]: http://www.hiltmon.com/blog/2012/02/20/the-markdown-mindset/
[disqus]: http://disqus.com/
[twitter]: https://twitter.com/PuerkitoBio
[crw]: http://www.compterenduweb.com/
