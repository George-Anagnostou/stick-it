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

      // Circular arrangement
      const symbolCount = card.length;
      const radius = 60;

      card.forEach((sticker, i) => {
        const img = document.createElement("img");
        img.src = `/uploads/${sticker}`;
        img.onerror = () => console.error(`Failed to load: ${img.src}`);
        img.setAttribute("src", img.src);

        // Position in a circle
        const angle = (i / symbolCount) * 2 * Math.PI;
        const x = radius * Math.cos(angle);
        const y = radius * Math.sin(angle);
        img.style.left = `calc(50% + ${x}px)`;
        img.style.top = `calc(50% + ${y}px)`;
        img.style.transform = "translate(-50%, -50%)";

        symbolsDiv.appendChild(img);
      });

      cardDiv.appendChild(symbolsDiv);
      deckDiv.appendChild(cardDiv);
    });

    document.getElementById("exportButton").style.display = "block";
  } catch (error) {
    console.error("Error:", error);
    alert("Something went wrong!");
  }
});

document.getElementById("exportButton").addEventListener("click", () => {
  window.location.href = "/export";
});
