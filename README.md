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
