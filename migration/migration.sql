USE forum_gaming;

-- Catégories
INSERT INTO categories (id, nom) VALUES
(1, 'General'),
(2, 'Technique'),
(3, 'Jeux Vidéo'),
(4, 'Entraide')
ON DUPLICATE KEY UPDATE nom=VALUES(nom);

-- Utilisateurs
-- Mot de passe en clair pour les comptes de démonstration : Password!123
-- (le hash ci-dessous est le SHA-512 de ce mot de passe, comme le génère l'API)
INSERT INTO utilisateurs (id, username, email, password_hash, role) VALUES
(1, 'admin', 'admin@example.com', '0b8ac2c4f042ead2a9f84038c2d692265517dfed6aec2d24010c2a2fb789fa603dedc404a511c184b2622142fa9922ab9c31107a31a114904abbf9d40bd802cc', 'admin'),
(2, 'alex_gaming', 'alex@example.com', '0b8ac2c4f042ead2a9f84038c2d692265517dfed6aec2d24010c2a2fb789fa603dedc404a511c184b2622142fa9922ab9c31107a31a114904abbf9d40bd802cc', 'user'),
(3, 'sophie_dev', 'sophie@example.com', '0b8ac2c4f042ead2a9f84038c2d692265517dfed6aec2d24010c2a2fb789fa603dedc404a511c184b2622142fa9922ab9c31107a31a114904abbf9d40bd802cc', 'user'),
(4, 'ludo_retro', 'ludo@example.com', '0b8ac2c4f042ead2a9f84038c2d692265517dfed6aec2d24010c2a2fb789fa603dedc404a511c184b2622142fa9922ab9c31107a31a114904abbf9d40bd802cc', 'user'),
(5, 'chloe_fps', 'chloe@example.com', '0b8ac2c4f042ead2a9f84038c2d692265517dfed6aec2d24010c2a2fb789fa603dedc404a511c184b2622142fa9922ab9c31107a31a114904abbf9d40bd802cc', 'user'),
(6, 'max_hardware', 'max@example.com', '0b8ac2c4f042ead2a9f84038c2d692265517dfed6aec2d24010c2a2fb789fa603dedc404a511c184b2622142fa9922ab9c31107a31a114904abbf9d40bd802cc', 'user')
ON DUPLICATE KEY UPDATE 
    username=VALUES(username), 
    email=VALUES(email), 
    password_hash=VALUES(password_hash), 
    role=VALUES(role);

-- Fils de discussion (Topics)
INSERT INTO fils_de_discussion (id, titre, etat, auteur_id, created_at) VALUES
(1, 'Saison 2 de Crown of Embers — Vos avis et retours', 'ouvert', 1, '2026-06-10 10:00:00'),
(2, 'Elden Ring DLC : Partagez vos meilleurs builds !', 'ouvert', 2, '2026-06-11 11:30:00'),
(3, 'Nintendo Switch 2 : Quelles sont vos attentes ?', 'ouvert', 4, '2026-06-12 14:15:00'),
(4, 'Velocity GP : Avis sur le nouveau patch 1.1', 'ouvert', 5, '2026-06-13 16:00:00'),
(5, 'Votre avis sur Pixel Tactics ?', 'ouvert', 4, '2026-06-14 09:45:00'),
(6, 'GTA 6 : Vos théories après le dernier trailer', 'ouvert', 2, '2026-06-14 18:20:00'),
(7, 'Problème d\'écran noir au lancement de Crown of Embers', 'ouvert', 3, '2026-06-15 08:30:00'),
(8, 'Config PC pour faire tourner GTA 6 en Ultra', 'ouvert', 6, '2026-06-15 13:00:00'),
(9, 'Cherche de l\'aide pour le boss final de Crown of Embers', 'ouvert', 5, '2026-06-15 20:10:00'),
(10, 'Vos plus beaux screenshots de jeux', 'ouvert', 2, '2026-06-16 06:00:00')
ON DUPLICATE KEY UPDATE 
    titre=VALUES(titre), 
    etat=VALUES(etat), 
    auteur_id=VALUES(auteur_id),
    created_at=VALUES(created_at);

-- Catégories associées aux fils
INSERT INTO fils_categories (fil_id, categorie_id) VALUES
(1, 3), -- Crown of Embers -> Jeux Vidéo
(2, 3), -- Elden Ring DLC -> Jeux Vidéo
(3, 1), -- Switch 2 -> General
(4, 3), -- Velocity GP -> Jeux Vidéo
(5, 3), -- Pixel Tactics -> Jeux Vidéo
(6, 1), -- GTA 6 -> General
(7, 2), -- Écran noir -> Technique
(8, 2), -- Config PC -> Technique
(9, 4), -- Aide boss -> Entraide
(10, 1) -- Screenshots -> General
ON DUPLICATE KEY UPDATE fil_id=VALUES(fil_id), categorie_id=VALUES(categorie_id);

-- Messages
INSERT INTO messages (id, contenu, fil_id, auteur_id, score_popularite, created_at) VALUES
-- Thread 1
(1, 'Bienvenue sur le fil officiel pour discuter de la Saison 2 de Crown of Embers. Que pensez-vous du nouveau biome volcanique ?', 1, 1, 3, '2026-06-10 10:00:00'),
(2, 'Franchement, le biome volcanique est incroyable visuellement ! La lave a des effets de distorsion thermique super bien faits. Par contre, le boss de la zone est hyper dur !', 1, 2, 5, '2026-06-10 10:15:00'),
(3, 'Carrément d\'accord, le boss volcanique m\'a tué au moins 15 fois avant que je comprenne son pattern de phase 2. Mais le loot en vaut la peine, l\'épée de braise est surpuissante.', 1, 5, 2, '2026-06-10 10:30:00'),
(4, 'Ravi que la mise à jour vous plaise ! N\'hésitez pas à partager vos stratégies pour vaincre le boss de zone.', 1, 1, 1, '2026-06-10 11:00:00'),

-- Thread 2
(5, 'Salut tout le monde ! Je cherche des idées de builds intéressants pour le DLC d\'Elden Ring. J\'ai fini le jeu de base en Force/Foi, vous me conseillez quoi ?', 2, 2, 2, '2026-06-11 11:30:00'),
(6, 'Si tu veux t\'amuser, tente un build Dextérité/Esotérisme avec les doubles lames inversées. Le moveset est ultra fluide et les saignements s\'accumulent super vite.', 2, 5, 6, '2026-06-11 11:45:00'),
(7, 'Moi je joue un build pur Force avec le marteau géant en acier noir. C\'est lent, mais chaque coup brise la posture des boss en deux ou trois sauts. Un vrai régal !', 2, 6, 4, '2026-06-11 12:00:00'),
(8, 'Ah cool ! Je pense que je vais tester les doubles lames inversées, le côté ninja agile me tente bien pour ce DLC. Merci pour les conseils !', 2, 2, 1, '2026-06-11 12:30:00'),

-- Thread 3
(9, 'Hello ! On approche petit à petit des annonces officielles pour la Nintendo Switch 2. Qu\'est-ce que vous attendez le plus ? Rétrocompatibilité, meilleure autonomie, plus de puissance ?', 3, 4, 4, '2026-06-12 14:15:00'),
(10, 'Pour moi la rétrocompatibilité est obligatoire. Si je ne peux pas jouer à mes jeux Switch actuels avec de meilleures performances, je ne l\'achèterai pas day one.', 3, 2, 8, '2026-06-12 14:30:00'),
(11, 'Pareil, et j\'espère un écran OLED d\'office, pas de retour au LCD bas de gamme. Niveau puissance, si on peut avoir du 1080p 60fps stable sur tous les jeux, ce serait déjà top.', 3, 6, 5, '2026-06-12 14:45:00'),
(12, 'Je pense qu\'ils vont faire un effort sur le processeur, peut-être avec du DLSS de Nvidia pour monter en résolution sur la TV. Hâte d\'en voir plus en tout cas !', 3, 4, 3, '2026-06-12 15:00:00'),

-- Thread 4
(13, 'Le patch 1.1 de Velocity GP vient de sortir ! Est-ce que vous trouvez aussi que la physique des collisions a été améliorée ? Avant c\'était n\'importe quoi.', 4, 5, 1, '2026-06-13 16:00:00'),
(14, 'Oui, je viens de faire 3 courses en ligne et c\'est le jour et la nuit. Les voitures ne s\'envolent plus au moindre contact. C\'est beaucoup plus réaliste et compétitif.', 4, 3, 4, '2026-06-13 16:20:00'),
(15, 'Super nouvelle ! Et que dis-tu des nouveaux circuits urbains ? Celui de Tokyo est super technique avec toutes ses épingles.', 4, 5, 2, '2026-06-13 16:40:00'),
(16, 'Tokyo est génial ! C\'est serré, pas le droit à l\'erreur. J\'ai hâte de monter dans le classement ce soir.', 4, 3, 2, '2026-06-13 17:00:00'),

-- Thread 5
(17, 'Quelqu\'un a testé Pixel Tactics ? J\'hésite à le prendre sur Switch 2. Les critiques sont bonnes mais j\'ai peur que ce soit trop répétitif.', 5, 4, 1, '2026-06-14 09:45:00'),
(18, 'Si tu aimes les jeux de stratégie au tour par tour en pixel art, fonce ! C\'est très tactique, avec une grosse profondeur dans la synergie des unités. Pas répétitif du tout.', 5, 2, 5, '2026-06-14 10:10:00'),
(19, 'Merci ! Je cherchais justement un petit jeu indé sympa pour mes trajets en train.', 5, 4, 0, '2026-06-14 10:30:00'),

-- Thread 6
(20, 'Le dernier trailer de GTA 6 est incroyable. Le niveau de détails sur la plage de Vice City est fou. Quelles sont vos théories sur le scénario de Lucia et Jason ?', 6, 2, 4, '2026-06-14 18:20:00'),
(21, 'Je parie sur une histoire à la Bonnie & Clyde avec un système de confiance mutuelle qui évolue selon nos choix. Si tu trahis ton partenaire, la fin sera différente.', 6, 5, 7, '2026-06-14 18:45:00'),
(22, 'Ce serait une excellente mécanique ! Rockstar a déposé un brevet sur des comportements d\'IA avancés il y a quelque temps, ça pourrait coller avec ça.', 6, 1, 3, '2026-06-14 19:10:00'),

-- Thread 7
(23, 'Bonjour, j\'ai un problème technique sur PC : écran noir direct au lancement de Crown of Embers. Le son tourne en fond mais rien ne s\'affiche. Quelqu\'un a une solution ?', 7, 3, 1, '2026-06-15 08:30:00'),
(24, 'Salut ! Essaie de mettre à jour tes pilotes graphiques ou de désactiver l\'overlay Discord/Steam. Souvent c\'est ça qui bloque le rendu du jeu.', 7, 6, 6, '2026-06-15 08:50:00'),
(25, 'Merci ! C\'était bien l\'overlay Discord qui faisait planter le jeu. Désactivé, et maintenant ça tourne nikel !', 7, 3, 2, '2026-06-15 09:15:00'),

-- Thread 8
(26, 'Quelle configuration minimale/recommandée selon vous pour faire tourner GTA 6 en Ultra 1440p/4K sur PC ? J\'envisage de changer de config d\'ici la sortie.', 8, 6, 3, '2026-06-15 13:00:00'),
(27, 'Il faudra probablement au minimum une RTX 4070 Ti ou équivalent chez AMD, avec un bon processeur Ryzen 7 pour gérer l\'IA de la foule qui s\'annonce massive.', 8, 3, 5, '2026-06-15 13:30:00'),
(28, 'Ok, donc ma RTX 3070 va commencer à être un peu juste... Je vais commencer à économiser pour une RTX 5000 d\'ici là !', 8, 6, 2, '2026-06-15 14:00:00'),

-- Thread 9
(29, 'Je n\'arrive pas à passer la troisième phase du boss final dans Crown of Embers. Des astuces sur les éléments à utiliser ou les patterns ?', 9, 5, 0, '2026-06-15 20:10:00'),
(30, 'Utilise des sorts ou de l\'équipement résistant à la foudre pour cette phase. Il charge son attaque ultime à 50% de PV, c\'est le moment idéal pour passer derrière lui et bourrer.', 9, 2, 4, '2026-06-15 20:30:00'),
(31, 'Génial, j\'ai appliqué le conseil et j\'ai enfin réussi à le battre ! L\'esquive vers l\'arrière au moment du saut était la clé.', 9, 5, 1, '2026-06-15 20:50:00'),

-- Thread 10
(32, 'Partagez ici vos plus beaux screenshots pris en jeu ! Je commence avec ce magnifique coucher de soleil dans Cyberpunk.', 10, 2, 8, '2026-06-16 06:00:00'),
(33, 'Magnifique ton screen ! De mon côté, voici une capture de Zelda sur Switch, la direction artistique compense largement la technique.', 10, 4, 5, '2026-06-16 06:30:00'),
(34, 'La DA de Zelda est intemporelle. Très beau cliché !', 10, 2, 3, '2026-06-16 07:00:00')
ON DUPLICATE KEY UPDATE 
    contenu=VALUES(contenu), 
    fil_id=VALUES(fil_id), 
    auteur_id=VALUES(auteur_id), 
    score_popularite=VALUES(score_popularite),
    created_at=VALUES(created_at);

-- Réactions
INSERT INTO reactions (utilisateur_id, message_id, type) VALUES
-- Reactions on message 2
(1, 2, 'like'),
(3, 2, 'like'),
(4, 2, 'like'),
(5, 2, 'like'),
(6, 2, 'like'),
-- Reactions on message 6
(1, 6, 'like'),
(2, 6, 'like'),
(3, 6, 'like'),
(4, 6, 'like'),
(6, 6, 'like'),
-- Reactions on message 10
(1, 10, 'like'),
(3, 10, 'like'),
(4, 10, 'like'),
(5, 10, 'like'),
(6, 10, 'like'),
(2, 10, 'like'),
-- Reactions on message 21
(1, 21, 'like'),
(2, 21, 'like'),
(3, 21, 'like'),
(4, 21, 'like'),
(6, 21, 'like'),
-- Reactions on message 32
(1, 32, 'like'),
(3, 32, 'like'),
(4, 32, 'like'),
(5, 32, 'like'),
(6, 32, 'like'),
(2, 32, 'like')
ON DUPLICATE KEY UPDATE type=VALUES(type);
