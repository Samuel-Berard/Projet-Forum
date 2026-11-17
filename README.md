# Projet Forum — Plateforme communautaire « Forum Gaming »

Plateforme de discussion communautaire sur le thème du **jeu vidéo**, inspirée de JeuxVideos.com. Les utilisateurs peuvent s'inscrire, ouvrir des fils de discussion, publier des
messages, réagir (like / dislike), filtrer par catégorie, rechercher des sujets, et les
administrateurs disposent d'un tableau de bord de modération.

---


## Membres de l'équipe

| Nom | Rôle |
|-----|------|
| **Clément Belmondo** | Développement |
| **Samuel Berard** | Développement |

---

## Technologies utilisées

| Domaine | Technologie |
|---------|-------------|
| Langage backend | **Go**|
| Routeur HTTP | `gorilla/mux` |
| Authentification | JWT — `golang-jwt/jwt/v5` |
| Base de données | **MySQL**|
| Hachage des mots de passe | **SHA-512** |

---

## Architecture

Le projet est un **client serveur organisé en couches** , regroupé dans un
seul dépôt et un seul module Go. Il est découpé en deux serveurs HTTP :


Chaque serveur respecte une séparation claire des responsabilités :

- **Router** — définition des routes.
- **Controller** — réception des requêtes, validation des entrées, appel de la logique métier.
- **Service** — logique métier.
- **Repository** — accès aux données et communication avec MySQL.
- **Model / DTO** — représentation des données manipulées.
- **View** — affichage des pages HTML (côté client).

---



## Prérequis

Avant d'installer le projet, assurez-vous d'avoir :

- **Go** 1.26 ou supérieur 
- **MySQL** 
- **Git** — pour cloner le dépôt

Vérification rapide :

```bash
go version
mysql --version
```

---

## Installation

### 1. Cloner le dépôt

```bash
git clone https://github.com/Samuel-Berard/Projet-Forum.git
cd Projet-Forum
```

### 2. Créer la base de données et insérer les données de test

Le script `script.sql` crée la base `forum_gaming` et ses tables.
Le script `migration.sql` insère les catégories, utilisateurs, fils et messages de démonstration.


### 3. Configurer les variables d'environnement

Créez un fichier `.env` à la **racine** du projet :

```dotenv
DB_USER=root
DB_PWD=votre_mot_de_passe
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=forum_gaming

# Optionnel — URL de l'API consommée par le client
# (valeur par défaut si absente)
API_BASE_URL=http://localhost:8080/api
```

### 4. Installer les dépendances Go

```bash
go mod download
```

---

## Lancement de l'application

L'application nécessite **deux serveurs** lancés en parallèle (dans deux terminaux).

**Terminal 1 — API (backend, port 8080) :**

```bash
go run ./api
```

**Terminal 2 — Client (frontend, port 3000) :**

```bash
go run ./client
```

Ouvrez ensuite votre navigateur sur :

> **http://localhost:3000**

---

## Comptes de démonstration

Les données de test (`migration.sql`) créent les comptes suivants.
Mot de passe commun : **`Password!123`**

| Identifiant | Email | Rôle |
|-------------|-------|------|
| `admin` | admin@example.com | **administrateur** |
| `alex_gaming` | alex@example.com | utilisateur |
| `sophie_dev` | sophie@example.com | utilisateur |
| `ludo_retro` | ludo@example.com | utilisateur |
| `chloe_fps` | chloe@example.com | utilisateur |
| `max_hardware` | max@example.com | utilisateur |

---

## Structure du projet

```
Projet-Forum/
├── api/                  # Serveur API (backend, :8080)
│   ├── main.go
│   ├── app/              # Assemblage de l'application (injection de dépendances)
│   ├── config/           # Chargement .env et connexion MySQL
│   ├── routers/          # Définition des routes /api/...
│   ├── controllers/      # Réception des requêtes, validation
│   ├── services/         # Logique métier
│   ├── repositories/     # Accès aux données (SQL)
│   ├── models/           # Entités
│   ├── dto/              # Objets de transfert de données
│   ├── middleware/       # Authentification JWT, contrôle des bannis
│   ├── auth/             # Génération et validation des JWT
│   ├── helper/ utils/    # Fonctions utilitaires (dont hachage SHA-512)
│   └── ...
├── client/               # Serveur web (frontend, :3000)
│   ├── main.go
│   ├── app/              # Assemblage de l'application cliente
│   ├── routers/          # Routes des pages + assets
│   ├── controllers/      # Rendu des pages
│   ├── services/         # Appels métier
│   ├── api/              # Client HTTP vers l'API backend
│   ├── templates/        # Pages HTML (index, forum, thread, login, admin...)
│   └── config/
├── static/               # CSS, JS et images servis sur /static
│   ├── css/
│   ├── js/
│   └── img/
├── uploads/              # Images envoyées par les utilisateurs (bonus)
├── migration/
│   ├── script.sql        # Création de la base et des tables
│   └── migration.sql     # Données de test
├── go.mod / go.sum
└── README.md
```

### Schéma de la base de données

`utilisateurs`, `categories`, `fils_de_discussion`, `fils_categories` (liaison N–N),
`messages`, `reactions`. Les suppressions sont propagées via `ON DELETE CASCADE`
(supprimer un fil supprime ses messages et leurs réactions).

## Exploitation et reprise sur incident

Les commandes ci-dessous utilisent les noms Scalingo à remplacer par ceux de
l'environnement concerné. Le projet utilise **MySQL** : les procédures
PostgreSQL ne s'appliquent pas à cette base. Ne jamais afficher ni copier les
valeurs des variables contenant des identifiants ou des mots de passe.

### Diagnostiquer une panne avec les logs

```bash
# Logs récents de l'application
scalingo --app <app-api> logs --lines 200

# Retrouver l'identifiant de l'addon MySQL (ad-...)
scalingo --app <app-api> addons

# Logs récents de la base
scalingo --app <app-api> logs --addon <id-addon-mysql> --lines 200
```

Trier les lignes par leur horodatage, puis vérifier les erreurs de démarrage,
la connexion à MySQL, les réponses HTTP et les événements de déploiement.
Pour exercer le diagnostic sans interrompre le service public, utiliser une
application de staging : y configurer temporairement un `DB_HOST` inaccessible,
puis redémarrer son conteneur web et lire les logs :

```bash
# Sauvegarder la valeur de DB_HOST hors du dépôt, puis provoquer l'échec
scalingo --app <app-staging> env-set DB_HOST=127.0.0.1
scalingo --app <app-staging> ps
scalingo --app <app-staging> restart web-1
scalingo --app <app-staging> logs --lines 200

# Restaurer la valeur précédente de DB_HOST et vérifier le redémarrage
scalingo --app <app-staging> env-set DB_HOST=<ancienne-valeur>
scalingo --app <app-staging> restart web-1
scalingo --app <app-staging> logs --lines 200
```

L'échec attendu est signalé au démarrage par `Erreur ping base de donnees`
(voir `api/config/db.config.go`). **Ne pas provoquer cette panne en production.**

### Revenir à un déploiement précédent

Identifier le dernier commit fonctionnel dans l'historique des déploiements,
puis vérifier que la migration de base de données du commit fautif n'empêche
pas l'ancienne version de fonctionner. Pour une application liée au remote
Git `scalingo` et déployée sur `master`, le redéploiement du commit connu se
fait ainsi :

```bash
git push scalingo <commit-fonctionnel>:master --force
```

Cette commande remplace la branche distante et déclenche un déploiement.
Vérifier le résultat dans les logs et tester le site après le déploiement.
Adapter le nom du remote et de la branche si l'application n'utilise pas
`scalingo`/`master`. Un rollback du code ne restaure pas la base de données.

### Sauvegarder et restaurer MySQL

Les sauvegardes automatiques dépendent du plan de l'addon. Vérifier leur
présence et leur rétention dans l'onglet **Backups** du tableau de bord. Une
sauvegarde manuelle peut être déclenchée et vérifiée avec :

```bash
scalingo --app <app-api> --addon <id-addon-mysql> backups-create
scalingo --app <app-api> --addon <id-addon-mysql> backups
scalingo --app <app-api> --addon <id-addon-mysql> backups-download --backup <id-backup>
```

Conserver le fichier téléchargé dans un emplacement protégé. La restauration
remplace les données et entraîne une interruption de service : la tester
d'abord sur un environnement de staging et confirmer l'identifiant du backup
avant toute restauration en production. Depuis un poste disposant du client
MySQL, télécharger et extraire l'archive, ouvrir le tunnel vers l'URL de
connexion MySQL configurée dans Scalingo, puis importer le dump :

```bash
tar -xzf <archive-backup.tar.gz>
scalingo --app <app-api> db-tunnel <variable-url-mysql>
mysql --user=<utilisateur> --password --host=127.0.0.1 --port=<port-tunnel> <base> < <fichier-dump.sql>
```

`--password` demande le mot de passe de façon interactive. Utiliser les
paramètres de connexion de l'environnement cible sans les inscrire dans le
dépôt. Le format de l'archive peut contenir un dump dont le nom diffère de
`<fichier-dump.sql>` : vérifier le contenu extrait avant l'import. Après
restauration, contrôler les tables et effectuer un test fonctionnel de
l'application.

### Vérifier les ressources et le coût

Dans **Metrics** du tableau de bord Scalingo, relever régulièrement la mémoire
et le CPU, notamment aux heures de pointe. Confirmer que le conteneur de taille
S ne s'approche pas durablement de ses limites avant de conclure qu'il suffit.
Comparer les montants réellement facturés (application et addon MySQL) à
l'estimation mensuelle et noter la période mesurée :

| Mesure mensuelle | Estimation | Réel | Période / source |
|---|---:|---:|---|
| Application et conteneurs | À renseigner | À relever | Facturation Scalingo |
| Addon MySQL et sauvegardes | À renseigner | À relever | Facturation Scalingo |
| Total | À renseigner | À relever | Même période |

**Validation opérationnelle à compléter :** les commandes et procédures
ci-dessus sont un runbook, pas un compte rendu d'exécution. Après les avoir
exécutées sur le compte Scalingo du projet (en particulier l'erreur contrôlée,
le rollback et une restauration de test), remplacer les champs « À relever »,
consigner la date et le résultat, et ne déclarer testée que la restauration
effectivement vérifiée.
