// thread.js — comportements de la page d'un fil de discussion.
// (Déplacé depuis le <script> et les attributs onclick/onsubmit/onchange du HTML.)

document.addEventListener("DOMContentLoaded", () => {
  const smileys = [
    { code: ":)", url: "https://image.jeuxvideo.com/smileys_img/1.gif" },
    { code: ":snif:", url: "https://image.jeuxvideo.com/smileys_img/20.gif" },
    { code: ":gba:", url: "https://image.jeuxvideo.com/smileys_img/17.gif" },
    { code: ":g)", url: "https://image.jeuxvideo.com/smileys_img/3.gif" },
    { code: ":-)", url: "https://image.jeuxvideo.com/smileys_img/46.gif" },
    { code: ":hap:", url: "https://image.jeuxvideo.com/smileys_img/18.gif" },
    { code: ":oui:", url: "https://image.jeuxvideo.com/smileys_img/37.gif" },
    { code: ":(", url: "https://image.jeuxvideo.com/smileys_img/45.gif" },
    { code: ":cool:", url: "https://image.jeuxvideo.com/smileys_img/26.gif" },
    { code: ":rire:", url: "https://image.jeuxvideo.com/smileys_img/39.gif" },
    { code: ":coeur:", url: "https://image.jeuxvideo.com/smileys_img/54.gif" },
    { code: ":merci:", url: "https://image.jeuxvideo.com/smileys_img/58.gif" }
  ];

  // Rendu "BBcode" + smileys du contenu de chaque message.
  document.querySelectorAll(".msg-body").forEach(el => {
    let text = el.textContent;

    text = text.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
    text = text.replace(/'''(.*?)'''/g, "<strong>$1</strong>");
    text = text.replace(/''(.*?)''/g, "<em>$1</em>");
    text = text.replace(/&lt;code&gt;(.*?)&lt;\/code&gt;/g, '<div class="code-preview">$1</div>');
    text = text.replace(/&lt;spoil&gt;(.*?)&lt;\/spoil&gt;/g, '<span class="spoil-preview"><button class="btn-spoil-red" type="button">SPOIL</button><span class="spoil-text">$1</span></span>');

    smileys.forEach(s => {
      const escapedCode = s.code.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
      text = text.replace(new RegExp(escapedCode, "g"), `<img src="${s.url}" alt="${s.code}">`);
    });

    text = text.replace(/!\[image\]\((.*?)\)/g, '<img src="$1" alt="image" class="msg-img">');
    text = text.replace(/\[video\]\((.*?)\)/g, '<div class="video-embed"><i class="fa-solid fa-play"></i> Vidéo : <a href="$1" target="_blank">$1</a></div>');
    text = text.replace(/\n/g, "<br>");

    el.innerHTML = text;
  });

  // Révéler/masquer un spoiler au clic.
  document.addEventListener("click", (e) => {
    if (e.target.classList.contains("btn-spoil-red")) {
      const t = e.target.nextElementSibling;
      if (t && t.classList.contains("spoil-text")) {
        t.classList.toggle("revealed");
      }
    }
  });

  // --- Handlers déplacés depuis les attributs HTML (onclick / onsubmit / onchange) ---

  // Modifier le titre du fil
  document.querySelectorAll('[data-action="toggle-thread-title"]').forEach(btn => {
    btn.addEventListener("click", toggleEditThreadTitle);
  });

  // Éditer un message (le bouton porte data-msg-id)
  document.querySelectorAll('[data-action="toggle-edit-msg"]').forEach(btn => {
    btn.addEventListener("click", () => toggleEditMessage(btn.dataset.msgId));
  });

  // Réagir (like / dislike) : data-msg-id + data-reaction
  document.querySelectorAll('[data-action="react"]').forEach(btn => {
    btn.addEventListener("click", () => sendReaction(btn.dataset.msgId, btn.dataset.reaction));
  });

  // Le select d'état envoie le formulaire au changement
  document.querySelectorAll("[data-autosubmit]").forEach(sel => {
    sel.addEventListener("change", () => sel.form.submit());
  });

  // Demande de confirmation avant l'envoi (formulaires avec data-confirm)
  document.querySelectorAll("form[data-confirm]").forEach(form => {
    form.addEventListener("submit", (e) => {
      if (!confirm(form.dataset.confirm)) {
        e.preventDefault();
      }
    });
  });
});

// Affiche/cache le formulaire de modification du titre.
function toggleEditThreadTitle() {
  document.getElementById("edit-thread-title-box")?.classList.toggle("open");
}

// Affiche/cache l'éditeur d'un message (et masque/réaffiche son texte).
function toggleEditMessage(msgId) {
  document.getElementById("edit-box-" + msgId)?.classList.toggle("open");
  document.getElementById("msg-body-" + msgId)?.classList.toggle("hidden");
}

// Envoie une réaction à l'API et met à jour le score sans recharger la page.
async function sendReaction(msgId, reactionType) {
  try {
    const response = await fetch(`/messages/${msgId}/react?type=${reactionType}`, { method: "POST" });
    const data = await response.json();
    if (!response.ok) {
      alert(data.erreur || "Une erreur est survenue");
      return;
    }

    const scoreEl = document.getElementById("score-" + msgId);
    let currentScore = parseInt(scoreEl.textContent.replace("+", ""));

    const likeBtn = scoreEl.previousElementSibling;
    const dislikeBtn = scoreEl.nextElementSibling;
    const hadLike = likeBtn.classList.contains("active");
    const hadDislike = dislikeBtn.classList.contains("active");

    if (reactionType === "like") {
      if (hadLike) return;
      likeBtn.classList.add("active");
      if (hadDislike) { dislikeBtn.classList.remove("active"); currentScore += 2; }
      else { currentScore += 1; }
    } else {
      if (hadDislike) return;
      dislikeBtn.classList.add("active");
      if (hadLike) { likeBtn.classList.remove("active"); currentScore -= 2; }
      else { currentScore -= 1; }
    }

    scoreEl.textContent = (currentScore > 0 ? "+" : "") + currentScore;
  } catch (err) {
    console.error(err);
    alert("Erreur lors de la réaction");
  }
}
