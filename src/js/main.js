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

            if (response.ok) {
                const result = await response.json();
                displayMetar(result);
            } else {
                const errorData = await response.json();
                output.textContent = `エラー: ${errorData.error}`;
            }
        } catch (error) {
            output.textContent = "通信エラーが発生しました。";
        }
    });
    
    function displayMetar(data) {
        let html = `<h3>METAR解析結果</h3>
                    <div class="metar-info">
                        <p><strong>空港コード:</strong> ${data.airport}</p>
                        <p><strong>観測時刻:</strong> ${data.time}</p>`;
        
        if (data.windDirection || data.windSpeed) {
            html += `<p><strong>風向/風速:</strong> ${data.windDirection || "N/A"} / ${data.windSpeed || "N/A"}</p>`;
        }
        
        if (data.visibility) {
            html += `<p><strong>視程:</strong> ${data.visibility}</p>`;
        }
        
        if (data.clouds && data.clouds.length > 0) {
            html += `<div class="cloud-info">
                        <p><strong>雲情報:</strong></p>
                        <ul>`;
            data.clouds.forEach((cloud, index) => {
                html += `<li>雲層${index + 1}: ${cloud.type} ${cloud.height}</li>`;
            });
            html += `</ul></div>`;
        }
        
        if (data.temperature) {
            html += `<p><strong>気温:</strong> ${data.temperature}</p>`;
        }
        
        if (data.dewPoint) {
            html += `<p><strong>露点温度:</strong> ${data.dewPoint}</p>`;
        }
        
        if (data.pressure) {
            html += `<p><strong>気圧:</strong> ${data.pressure}</p>`;
        }
        
        if (data.tempoInfo) {
            html += `<p><strong>一時的な天候:</strong> ${data.tempoInfo}</p>`;
        }
        
        if (data.remarks) {
            html += `<p><strong>備考:</strong> ${data.remarks}</p>`;
        }
        
        html += `</div>`;
        output.innerHTML = html;
    }
});