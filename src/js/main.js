document.addEventListener("DOMContentLoaded", () => {
    const form = document.getElementById("metar-form");
    const input = document.getElementById("metar-input");
    const output = document.getElementById("output");

    form.addEventListener("submit", async (event) => {
        event.preventDefault();
        
        const metarData = input.value.trim();
        if (!metarData) {
            output.textContent = "METARを入力してください。";
            return;
        }

        try {
            const response = await fetch("/api/metar", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ raw: metarData })
            });

            const result = await response.json();
            if (response.ok) {
                output.textContent = result.readable;
            } else {
                output.textContent = `エラー: ${result.error}`;
            }
        } catch (error) {
            output.textContent = "通信エラーが発生しました。";
        }
    });
});
