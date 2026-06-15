document.addEventListener("DOMContentLoaded", () => {
    const textarea = document.getElementById("message_topic");
    const previewBox = document.querySelector(".editor-preview-box");
    const formatButtons = document.querySelectorAll(".toolbar-btn");
    const togglePreviewBtn = document.querySelector(".toggle-switch");
    
    // Modal Smileys
    const smileyModal = document.getElementById("smiley-modal");
    const closeSmileyModal = document.getElementById("close-smiley-modal");
    const smileyGrid = document.getElementById("smiley-grid");

    const smileys = [
        { code: ":)", url: "https://image.jeuxvideo.com/smileys_img/1.gif" },
        { code: ":snif:", url: "https://image.jeuxvideo.com/smileys_img/20.gif" },
        { code: ":gba:", url: "https://image.jeuxvideo.com/smileys_img/17.gif" },
        { code: ":g)", url: "https://image.jeuxvideo.com/smileys_img/3.gif" },
        { code: ":-)", url: "https://image.jeuxvideo.com/smileys_img/46.gif" },
        { code: ":snif2:", url: "https://image.jeuxvideo.com/smileys_img/13.gif" },
        { code: ":bravo:", url: "https://image.jeuxvideo.com/smileys_img/69.gif" },
        { code: ":d)", url: "https://image.jeuxvideo.com/smileys_img/4.gif" },
        { code: ":hap:", url: "https://image.jeuxvideo.com/smileys_img/18.gif" },
        { code: ":ouch:", url: "https://image.jeuxvideo.com/smileys_img/22.gif" },
        { code: ":pacg:", url: "https://image.jeuxvideo.com/smileys_img/9.gif" },
        { code: ":cd:", url: "https://image.jeuxvideo.com/smileys_img/5.gif" },
        { code: ":-)))", url: "https://image.jeuxvideo.com/smileys_img/23.gif" },
        { code: ":ouch2:", url: "https://image.jeuxvideo.com/smileys_img/57.gif" },
        { code: ":pacd:", url: "https://image.jeuxvideo.com/smileys_img/10.gif" },
        { code: ":cute:", url: "https://image.jeuxvideo.com/smileys_img/nyu.gif" },
        { code: ":content:", url: "https://image.jeuxvideo.com/smileys_img/24.gif" },
        { code: ":p)", url: "https://image.jeuxvideo.com/smileys_img/7.gif" },
        { code: ":-p", url: "https://image.jeuxvideo.com/smileys_img/31.gif" },
        { code: ":noel:", url: "https://image.jeuxvideo.com/smileys_img/11.gif" },
        { code: ":oui:", url: "https://image.jeuxvideo.com/smileys_img/37.gif" },
        { code: ":(", url: "https://image.jeuxvideo.com/smileys_img/45.gif" },
        { code: ":peur:", url: "https://image.jeuxvideo.com/smileys_img/47.gif" },
        { code: ":question:", url: "https://image.jeuxvideo.com/smileys_img/2.gif" },
        { code: ":cool:", url: "https://image.jeuxvideo.com/smileys_img/26.gif" },
        { code: ":-(", url: "https://image.jeuxvideo.com/smileys_img/14.gif" },
        { code: ":coeur:", url: "https://image.jeuxvideo.com/smileys_img/54.gif" },
        { code: ":mort:", url: "https://image.jeuxvideo.com/smileys_img/21.gif" },
        { code: ":rire:", url: "https://image.jeuxvideo.com/smileys_img/39.gif" },
        { code: ":-((", url: "https://image.jeuxvideo.com/smileys_img/15.gif" },
        { code: ":fou:", url: "https://image.jeuxvideo.com/smileys_img/50.gif" },
        { code: ":sleep:", url: "https://image.jeuxvideo.com/smileys_img/27.gif" }
    ];

    if (!textarea) return;

    // Config des balises de formatage (Syntaxe JVC)
    const tags = {
        "Gras": { open: "'''", close: "'''" },
        "Italique": { open: "''", close: "''" },
        "Souligné": { open: "<u>", close: "</u>" },
        "Barré": { open: "<s>", close: "</s>" },
        "Liste": { open: "* ", close: "" },
        "Listes numérotées": { open: "# ", close: "" },
        "Citation": { open: "> ", close: "" },
        "Code": { open: "<code>", close: "</code>" },
        "Spoiler": { open: "<spoil>", close: "</spoil>" },
        "Image": { open: "![image](", close: ")" },
        "Vidéo": { open: "[video](", close: ")" }
    };

    // Insérer les balises au clic
    formatButtons.forEach(btn => {
        btn.addEventListener("click", () => {
            const title = btn.getAttribute("title");
            
            if (title === "Smileys") {
                if (smileyModal) smileyModal.style.display = "flex";
                return;
            }

            if (title === "Image") {
                const imageModal = document.getElementById("image-modal");
                if (imageModal) imageModal.style.display = "flex";
                return;
            }

            if (title === "Vidéo") {
                const videoModal = document.getElementById("video-modal");
                if (videoModal) videoModal.style.display = "flex";
                return;
            }

            const tag = tags[title];
            if (!tag) return;

            const start = textarea.selectionStart;
            const end = textarea.selectionEnd;
            const selectedText = textarea.value.substring(start, end);
            
            const textBefore = textarea.value.substring(0, start);
            const textAfter = textarea.value.substring(end, textarea.value.length);
            
            textarea.value = textBefore + tag.open + selectedText + tag.close + textAfter;
            
            // Placer le curseur au bon endroit
            if (selectedText.length === 0) {
                textarea.selectionStart = textarea.selectionEnd = start + tag.open.length;
            } else {
                textarea.selectionStart = textarea.selectionEnd = start + tag.open.length + selectedText.length + tag.close.length;
            }
            
            textarea.focus();
            updatePreview();
        });
    });

    // Mettre à jour la prévisualisation en direct
    const updatePreview = () => {
        if (!previewBox) return;
        let text = textarea.value;

        // Cacher/Afficher le disclaimer comme un watermark
        const disclaimer = document.querySelector(".editor-disclaimer");
        if (disclaimer) {
            if (text.length > 0) {
                disclaimer.style.opacity = "0";
                disclaimer.style.pointerEvents = "none";
            } else {
                disclaimer.style.opacity = "1";
                disclaimer.style.pointerEvents = "auto";
            }
        }
        
        // Échapper le HTML pour la sécurité avant de remplacer (très basique)
        text = text.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");

        // Rétablir les balises JVC après échappement
        text = text.replace(/&lt;u&gt;/g, "<u>").replace(/&lt;\/u&gt;/g, "</u>");
        text = text.replace(/&lt;s&gt;/g, "<s>").replace(/&lt;\/s&gt;/g, "</s>");
        text = text.replace(/&lt;code&gt;/g, "<code>").replace(/&lt;\/code&gt;/g, "</code>");
        text = text.replace(/&lt;spoil&gt;/g, "<spoil>").replace(/&lt;\/spoil&gt;/g, "</spoil>");

        // Remplacements JVC
        text = text.replace(/'''(.*?)'''/g, '<strong>$1</strong>');
        text = text.replace(/''(.*?)''/g, '<em>$1</em>');
        text = text.replace(/<code>(.*?)<\/code>/g, '<div class="code-preview">$1</div>');
        text = text.replace(/<spoil>(.*?)<\/spoil>/g, '<div class="spoil-preview"><button class="btn-spoil-red">SPOIL</button> <span class="spoil-text">Afficher</span></div>');
        
        // Remplacements Smileys
        smileys.forEach(s => {
            const escapedCode = s.code.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
            const regex = new RegExp(escapedCode, 'g');
            text = text.replace(regex, `<img src="${s.url}" alt="${s.code}">`);
        });

        // Remplacement Image URL
        text = text.replace(/!\[image\]\((.*?)\)/g, '<img src="$1" alt="image" style="max-width: 100%; max-height: 400px; border-radius: 4px; margin-top: 0.5rem;">');

        // Remplacement Vidéo URL
        text = text.replace(/\[video\]\((.*?)\)/g, '<div class="code-preview" style="background:#1d1e20;color:#4a8bed;display:inline-block;padding:0.5rem 1rem;margin-top:0.5rem;"><i class="fa-solid fa-play"></i> Vidéo : <a href="$1" target="_blank" style="color:#4a8bed;text-decoration:underline;">$1</a></div>');

        text = text.replace(/\n/g, '<br>');
        
        previewBox.innerHTML = text;
    };

    textarea.addEventListener("input", updatePreview);
    updatePreview();

    // Activer / Désactiver la prévisu avec le switch
    if (togglePreviewBtn && previewBox) {
        togglePreviewBtn.addEventListener("click", () => {
            togglePreviewBtn.classList.toggle("active");
            if (togglePreviewBtn.classList.contains("active")) {
                previewBox.style.display = "block";
                togglePreviewBtn.style.background = "var(--color-link)";
            } else {
                previewBox.style.display = "none";
                togglePreviewBtn.style.background = "#4a4c4f";
            }
        });
    }

    // Gérer la modal Smileys
    if (smileyGrid) {
        smileys.forEach(s => {
            const item = document.createElement("div");
            item.className = "smiley-item";

            const imgWrapper = document.createElement("div");
            imgWrapper.className = "smiley-img-wrapper";
            imgWrapper.innerHTML = `<img src="${s.url}" alt="${s.code}">`;
            
            const codeWrapper = document.createElement("div");
            codeWrapper.className = "smiley-code-wrapper";
            codeWrapper.textContent = s.code;

            item.appendChild(imgWrapper);
            item.appendChild(codeWrapper);

            const insertSmiley = () => {
                const start = textarea.selectionStart;
                const textBefore = textarea.value.substring(0, start);
                const textAfter = textarea.value.substring(textarea.selectionEnd, textarea.value.length);
                textarea.value = textBefore + " " + s.code + " " + textAfter;
                textarea.selectionStart = textarea.selectionEnd = start + s.code.length + 2;
                textarea.focus();
                updatePreview();
                smileyModal.style.display = "none";
            };

            item.addEventListener("click", insertSmiley);

            smileyGrid.appendChild(item);
        });
    }

    if (smileyModal && closeSmileyModal) {
        closeSmileyModal.addEventListener("click", () => {
            smileyModal.style.display = "none";
        });
        
        // Fermer en cliquant à l'extérieur
        smileyModal.addEventListener("click", (e) => {
            if (e.target === smileyModal) {
                smileyModal.style.display = "none";
            }
        });
    }

    // Gérer la modal Images (NoelShack)
    const imageModal = document.getElementById("image-modal");
    const closeImageModal = document.getElementById("close-image-modal");
    const btnDownloadUrl = document.getElementById("btn-download-url");
    const imageUrlInput = document.getElementById("image-url-input");
    const uploadDropzone = document.getElementById("upload-dropzone");

    if (imageModal && closeImageModal) {
        closeImageModal.addEventListener("click", () => {
            imageModal.style.display = "none";
        });
        
        imageModal.addEventListener("click", (e) => {
            if (e.target === imageModal) {
                imageModal.style.display = "none";
            }
        });
    }

    // Gérer l'ajout d'image par URL
    if (btnDownloadUrl && imageUrlInput) {
        btnDownloadUrl.addEventListener("click", () => {
            const url = imageUrlInput.value.trim();
            if (url) {
                const start = textarea.selectionStart;
                const textBefore = textarea.value.substring(0, start);
                const textAfter = textarea.value.substring(textarea.selectionEnd, textarea.value.length);
                const imgSyntax = `![image](${url})`;
                
                textarea.value = textBefore + imgSyntax + textAfter;
                textarea.selectionStart = textarea.selectionEnd = start + imgSyntax.length;
                textarea.focus();
                updatePreview();
                
                imageUrlInput.value = "";
                imageModal.style.display = "none";
            }
        });
    }

    // Simulation de l'upload local (drag and drop / click)
    if (uploadDropzone) {
        uploadDropzone.addEventListener("click", () => {
            alert("L'upload d'image vers le disque n'est pas encore branché au backend. Utilisez l'ajout par URL en attendant !");
        });
    }

    // Gérer la modal Vidéo
    const videoModal = document.getElementById("video-modal");
    const closeVideoModal = document.getElementById("close-video-modal");
    const btnAddVideo = document.getElementById("btn-add-video");
    const videoUrlInput = document.getElementById("video-url-input");

    if (videoModal && closeVideoModal) {
        closeVideoModal.addEventListener("click", () => {
            videoModal.style.display = "none";
        });
        
        videoModal.addEventListener("click", (e) => {
            if (e.target === videoModal) {
                videoModal.style.display = "none";
            }
        });
    }

    if (btnAddVideo && videoUrlInput) {
        btnAddVideo.addEventListener("click", () => {
            const url = videoUrlInput.value.trim();
            if (url) {
                const start = textarea.selectionStart;
                const textBefore = textarea.value.substring(0, start);
                const textAfter = textarea.value.substring(textarea.selectionEnd, textarea.value.length);
                const vidSyntax = `[video](${url})`;
                
                textarea.value = textBefore + vidSyntax + textAfter;
                textarea.selectionStart = textarea.selectionEnd = start + vidSyntax.length;
                textarea.focus();
                updatePreview();
                
                videoUrlInput.value = "";
                videoModal.style.display = "none";
            }
        });
    }
});
