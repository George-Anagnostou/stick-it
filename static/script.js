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

      const numberDiv = document.createElement("div");
      numberDiv.className = "card-number";
      numberDiv.textContent = `Card ${index + 1}`;
      cardDiv.appendChild(numberDiv);

      const symbolsDiv = document.createElement("div");
      symbolsDiv.className = "symbols";

      // Dynamically adjust grid based on symbol count
      const symbolCount = card.length;
      const columns = Math.ceil(Math.sqrt(symbolCount));
      symbolsDiv.style.gridTemplateColumns = `repeat(${columns}, 1fr)`;

      card.forEach((sticker) => {
        const img = document.createElement("img");
        img.src = `/uploads/${sticker}`;
        img.onerror = () => console.error(`Failed to load: ${img.src}`);
        img.setAttribute("src", img.src);
        symbolsDiv.appendChild(img);
      });

      cardDiv.appendChild(symbolsDiv);
      deckDiv.appendChild(cardDiv);
    });
  } catch (error) {
    console.error("Error:", error);
    alert("Something went wrong!");
  }
});
