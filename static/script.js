document.getElementById("uploadForm").addEventListener("submit", async (e) => {
  e.preventDefault();
  const formData = new FormData(e.target);

  try {
    const response = await fetch("/upload", {
      method: "POST",
      body: formData,
    });
    const result = await response.json();

    if (result.error) {
      alert(result.error);
      return;
    }

    console.log("Response:", result);

    const deckDiv = document.getElementById("deck");
    deckDiv.innerHTML = "";

    result.deck.forEach((card, index) => {
      const cardDiv = document.createElement("div");
      cardDiv.className = "card";
      cardDiv.innerHTML = `<strong>Card ${index + 1}</strong>`;

      card.forEach((sticker) => {
        console.log("Loading image:", `/uploads/${sticker}`);
        const img = document.createElement("img");
        img.src = `/uploads/${sticker}`;
        img.onload = () => console.log(`Loaded: ${img.src}`);
        img.onerror = () => console.error(`Failed to load: ${img.src}`);
        img.setAttribute("src", img.src);
        cardDiv.appendChild(img);

        //fetch(`/uploads/${sticker}`)
        //  .then((res) => res.blob())
        //  .then((blob) => {
        //    const img = document.createElement("img");
        //    const url = URL.createObjectURL(blob);
        //    img.src = url;
        //    img.onload = () => console.log(`Loaded via blob: ${sticker}`);
        //    img.onerror = () => console.error(`Failed via blob: ${sticker}`);
        //    cardDiv.appendChild(img);
        //  })
        //  .catch((err) => console.error("Blob fetch error:", err));
      });

      deckDiv.appendChild(cardDiv);
    });
  } catch (error) {
    console.error("Error:", error);
    alert("Something went wrong!");
  }
});
