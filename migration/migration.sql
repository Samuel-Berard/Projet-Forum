USE forum_gaming;

-- Catégories
INSERT INTO categories (nom) VALUES
('General'),
('Technique'),
('Jeux Vidéo'),
('Entraide');

-- Utilisateurs
-- Mot de passe en clair pour les comptes de démonstration : Password!123
-- (le hash ci-dessous est le SHA-512 de ce mot de passe, comme le génère l'API)
INSERT INTO utilisateurs (username, email, password_hash, role) VALUES
('admin', 'admin@example.com', '0b8ac2c4f042ead2a9f84038c2d692265517dfed6aec2d24010c2a2fb789fa603dedc404a511c184b2622142fa9922ab9c31107a31a114904abbf9d40bd802cc', 'admin'),
('alice', 'alice@example.com', '0b8ac2c4f042ead2a9f84038c2d692265517dfed6aec2d24010c2a2fb789fa603dedc404a511c184b2622142fa9922ab9c31107a31a114904abbf9d40bd802cc', 'user'),
('bob', 'bob@example.com', '0b8ac2c4f042ead2a9f84038c2d692265517dfed6aec2d24010c2a2fb789fa603dedc404a511c184b2622142fa9922ab9c31107a31a114904abbf9d40bd802cc', 'user');

-- Fils de discussion
INSERT INTO fils_de_discussion (titre, etat, auteur_id) VALUES
('Bienvenue sur le forum', 'ouvert', 1),
('Vos configs PC gaming', 'ouvert', 2),
('Entraide : problème de connexion', 'ouvert', 3);

-- Catégories associées aux fils
INSERT INTO fils_categories (fil_id, categorie_id) VALUES
(1, 1),
(2, 3),
(3, 4);

-- Messages
INSERT INTO messages (contenu, fil_id, auteur_id) VALUES
('Bienvenue à toutes et à tous sur le forum !', 1, 1),
('Partagez vos configurations PC ici.', 2, 2),
('Quelqu''un peut m''aider avec ma connexion ?', 3, 3);

-- Réactions
INSERT INTO reactions (utilisateur_id, message_id, type) VALUES
(2, 1, 'like'),
(3, 1, 'like');

-- Mise à jour du score de popularité en fonction des réactions
UPDATE messages SET score_popularite = 2 WHERE id = 1;
