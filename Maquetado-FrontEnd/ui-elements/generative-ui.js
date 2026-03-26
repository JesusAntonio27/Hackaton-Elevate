(function () {
  function escapeHtml(value) {
    return String(value ?? "")
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;")
      .replace(/'/g, "&#39;");
  }

  function markdownBasicoToHtml(md) {
    var raw = escapeHtml(md || "").replace(/\r\n/g, "\n");
    var lines = raw.split("\n");
    var out = [];
    var inList = false;
    var inTable = false;

    function inlineFormat(value) {
      return value
        .replace(/\*\*(.*?)\*\*/g, "<strong>$1</strong>")
        .replace(/`([^`]+)`/g, "<code>$1</code>")
        .replace(/(^|[^\*])\*(?!\s)([^*]+?)\*(?!\*)/g, "$1<em>$2</em>");
    }

    function closeList() {
      if (inList) {
        out.push("</ul>");
        inList = false;
      }
    }

    function closeTable() {
      if (inTable) {
        out.push("</tbody></table>");
        inTable = false;
      }
    }

    function parseTableCells(line) {
      return line
        .trim()
        .replace(/^\|/, "")
        .replace(/\|$/, "")
        .split("|")
        .map(function (cell) { return inlineFormat(cell.trim()); });
    }

    for (var i = 0; i < lines.length; i += 1) {
      var line = lines[i];
      var trimmed = line.trim();

      // Separador horizontal.
      if (/^---+$/.test(trimmed)) {
        closeList();
        closeTable();
        out.push("<hr>");
        continue;
      }

      // Headers # a ####
      var h = /^(#{1,4})\s+(.*)$/.exec(trimmed);
      if (h) {
        closeList();
        closeTable();
        var level = h[1].length + 1; // h2-h5 para no romper jerarquia visual.
        out.push("<h" + level + ">" + inlineFormat(h[2]) + "</h" + level + ">");
        continue;
      }

      // Listas - item o * item
      if (/^[-*]\s+/.test(trimmed)) {
        closeTable();
        if (!inList) {
          out.push("<ul>");
          inList = true;
        }
        out.push("<li>" + inlineFormat(trimmed.replace(/^[-*]\s+/, "")) + "</li>");
        continue;
      }
      closeList();

      // Tablas markdown simples
      var next = i + 1 < lines.length ? lines[i + 1].trim() : "";
      if (trimmed.includes("|") && /^\|?[\s:-]+\|[\s|:-]*$/.test(next)) {
        closeTable();
        var headers = parseTableCells(trimmed);
        out.push("<table><thead><tr>" + headers.map(function (c) { return "<th>" + c + "</th>"; }).join("") + "</tr></thead><tbody>");
        inTable = true;
        i += 1; // Saltar linea separadora de encabezado.
        continue;
      }
      if (inTable && trimmed.includes("|")) {
        var cells = parseTableCells(trimmed);
        out.push("<tr>" + cells.map(function (c) { return "<td>" + c + "</td>"; }).join("") + "</tr>");
        continue;
      }
      closeTable();

      if (trimmed.length === 0) {
        out.push("<br>");
      } else {
        out.push("<div>" + inlineFormat(trimmed) + "</div>");
      }
    }

    closeList();
    closeTable();
    return out.join("");
  }

  function horaActual() {
    return new Date().toLocaleTimeString("es-MX", { hour: "2-digit", minute: "2-digit" });
  }

  function renderTags(etiquetas) {
    if (!Array.isArray(etiquetas) || etiquetas.length === 0) return "";
    return etiquetas
      .map(function (tag) {
        var color = tag && tag.color ? escapeHtml(tag.color) : "blue";
        var texto = tag && tag.texto ? escapeHtml(tag.texto) : "";
        return '<span class="tag ' + color + '">' + texto + "</span>";
      })
      .join("");
  }

  function renderParetoInsight(uiData) {
    var filas = Array.isArray(uiData && uiData.filas) ? uiData.filas : [];
    var series = Array.isArray(uiData && uiData.series) ? uiData.series : [];
    var analisis = uiData && uiData.analisis ? markdownBasicoToHtml(uiData.analisis) : "";
    var chartSeries = series.length > 0 ? series.slice(0, 5) : filas.slice(0, 5).map(function (fila, idx) {
      var val = Number(String((fila && fila.valor) || "").replace(/[^\d.]/g, ""));
      return {
        etiqueta: fila && fila.etiqueta ? fila.etiqueta : "Causa " + (idx + 1),
        porcentaje: Number.isFinite(val) && val > 0 ? Math.min(val, 100) : Math.max(100 - idx * 15, 20),
        menciones: fila && fila.valor ? fila.valor : ""
      };
    });

    var chartHtml =
      '<div class="pareto-mini-chart">' +
      chartSeries
        .map(function (item, idx) {
          var pct = Number(item.porcentaje || 0);
          if (!Number.isFinite(pct)) pct = 0;
          pct = Math.max(0, Math.min(100, pct));
          return (
            '<div class="pareto-mini-row">' +
            '<div class="pareto-mini-label">' + escapeHtml(item.etiqueta || ("Causa " + (idx + 1))) + "</div>" +
            '<div class="pareto-mini-bar-bg"><div class="pareto-mini-bar' + (idx < 2 ? " top" : "") + '" style="width:' + pct + '%"></div></div>' +
            '<div class="pareto-mini-val">' + pct + "%</div>" +
            "</div>"
          );
        })
        .join("") +
      "</div>";

    return (
      '<div class="insight-card">' +
      '<div class="insight-label">' + escapeHtml(uiData && uiData.titulo_tarjeta ? uiData.titulo_tarjeta : "Insight") + "</div>" +
      chartHtml +
      (filas
        .map(function (fila) {
          return (
            '<div class="insight-row">' +
            '<span class="insight-key">' + escapeHtml(fila && fila.etiqueta) + "</span>" +
            '<span class="insight-val">' + escapeHtml(fila && fila.valor) + "</span>" +
            "</div>"
          );
        })
        .join("")) +
      (analisis ? '<div class="pareto-mini-analysis">' + analisis + "</div>" : "") +
      "</div>"
    );
  }

  function renderGeneric(uiData) {
    if (!uiData || !uiData.descripcion) return "";
    return (
      '<div class="insight-card">' +
      '<div class="insight-label">Detalle</div>' +
      '<div style="font-size:13px;line-height:1.5;">' + markdownBasicoToHtml(uiData.descripcion) + "</div>" +
      "</div>"
    );
  }

  var uiRenderers = {
    pareto_insight: renderParetoInsight
  };

  function normalizarRespuestaIA(payload) {
    if (!payload) return null;
    if (payload.respuesta_ia) return payload.respuesta_ia;
    if (payload.resultado) return payload.resultado;
    return payload;
  }

  function renderizarRespuestaIA(jsonIAResponse) {
    var respuesta = normalizarRespuestaIA(jsonIAResponse) || {};
    var msgs = document.getElementById("messages");
    if (!msgs) return;

    var uiType = respuesta.ui_type;
    var uiData = respuesta.ui_data || {};
    var renderer = uiRenderers[uiType] || renderGeneric;
    var bloqueUI = renderer(uiData);

    var botMsg = document.createElement("div");
    botMsg.className = "msg bot";
    botMsg.innerHTML =
      '<div class="msg-avatar bot-av">' +
      '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>' +
      "</div>" +
      "<div>" +
      '<div class="msg-bubble">' +
      markdownBasicoToHtml(respuesta.mensaje_principal || "No se recibio contenido de respuesta.") +
      (bloqueUI ? "<br><br>" + bloqueUI : "") +
      "</div>" +
      '<div style="margin-top:6px;">' + renderTags(uiData.etiquetas_tags) + "</div>" +
      '<div class="msg-time">' + horaActual() + "</div>" +
      "</div>";

    msgs.appendChild(botMsg);
    msgs.scrollTop = msgs.scrollHeight;
  }

  function actualizarFuentesRAG(fuentes) {
    var panelCuerpo = document.querySelector(".rp-body");
    if (!panelCuerpo) return;
    panelCuerpo.innerHTML = "";

    var lista = Array.isArray(fuentes) ? fuentes : [];
    if (lista.length === 0) {
      panelCuerpo.innerHTML = '<div class="source-card"><div class="source-text">Sin fuentes para esta consulta.</div></div>';
      return;
    }

    lista.forEach(function (fuente) {
      var score = Number(fuente.nps_score || 0);
      var colorVar = score <= 6 ? "rojo" : score <= 8 ? "naranja" : "verde";
      var card = document.createElement("div");
      card.className = "source-card";
      card.innerHTML =
        '<div class="source-top">' +
        '<span class="source-id">ID: ' + escapeHtml(fuente.id_respuesta || "N/A") + "</span>" +
        '<span class="source-score" style="color:var(--' + colorVar + ')">Score: ' + escapeHtml(score) + " · " + escapeHtml(fuente.clasificacion_nps || "N/A") + "</span>" +
        "</div>" +
        '<div class="source-text">"' + escapeHtml(fuente.comentario || "") + '"</div>' +
        '<div class="source-meta">Agente: ' + escapeHtml(fuente.agente_soporte || "N/A") + " · Área: " + escapeHtml(fuente.area_colaborador || "N/A") + "</div>";
      panelCuerpo.appendChild(card);
    });
  }

  function renderErrorBot(msg) {
    renderizarRespuestaIA({
      mensaje_principal: msg || "No pude procesar la consulta en este momento.",
      ui_type: "generic",
      ui_data: {
        descripcion: "Intenta nuevamente en unos segundos.",
        etiquetas_tags: [{ color: "red", texto: "Error de servicio" }]
      }
    });
  }

  window.ChatbotUI = {
    renderizarRespuestaIA: renderizarRespuestaIA,
    actualizarFuentesRAG: actualizarFuentesRAG,
    normalizarRespuestaIA: normalizarRespuestaIA,
    renderErrorBot: renderErrorBot
  };
})();
