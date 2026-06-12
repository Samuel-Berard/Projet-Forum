CREATE TABLE IF NOT EXISTS utilisateurs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role ENUM('user', 'admin') DEFAULT 'user',
    banned BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nom VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS fils_de_discussion (
    id INT AUTO_INCREMENT PRIMARY KEY,
    titre VARCHAR(255) NOT NULL,
    etat ENUM('ouvert', 'ferme', 'archive') DEFAULT 'ouvert',
    auteur_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (auteur_id) REFERENCES utilisateurs(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS fils_categories (
    fil_id INT NOT NULL,
    categorie_id INT NOT NULL,
    PRIMARY KEY (fil_id, categorie_id),
    FOREIGN KEY (fil_id) REFERENCES fils_de_discussion(id) ON DELETE CASCADE,
    FOREIGN KEY (categorie_id) REFERENCES categories(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    contenu TEXT NOT NULL,
    fil_id INT NOT NULL,
    auteur_id INT NOT NULL,
    score_popularite INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (fil_id) REFERENCES fils_de_discussion(id) ON DELETE CASCADE,
    FOREIGN KEY (auteur_id) REFERENCES utilisateurs(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS reactions (
    utilisateur_id INT NOT NULL,
    message_id INT NOT NULL,
    type ENUM('like', 'dislike') NOT NULL,
    PRIMARY KEY (utilisateur_id, message_id),
    FOREIGN KEY (utilisateur_id) REFERENCES utilisateurs(id) ON DELETE CASCADE,
    FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE
);

INSERT IGNORE INTO categories (nom) VALUES ('General'), ('Technique'), ('Jeux Vidéo'), ('Entraide');
