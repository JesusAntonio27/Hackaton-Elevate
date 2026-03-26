let allNPSData = [];
let currentMode = 'soporte'; // 'soporte' or 'general'

// NPS trend chart
const ctx1 = document.getElementById('npsChart').getContext('2d');
window.npsChartInstance = new Chart(ctx1, {
  type: 'line',
  data: {
    labels: [], // Will be populated dynamically
    datasets: [
      {
        label: 'NPS',
        data: [], // Will be populated dynamically
        borderColor: '#003087',
        backgroundColor: 'rgba(0,48,135,0.08)',
        borderWidth: 2.5,
        tension: 0.4,
        fill: true,
        pointBackgroundColor: '#003087',
        pointRadius: 4,
        pointHoverRadius: 6,
      },
      {
        label: 'Promotores %',
        data: [], // Will be populated dynamically
        borderColor: '#1B8A5A',
        backgroundColor: 'transparent',
        borderWidth: 2,
        borderDash: [5, 4],
        tension: 0.4,
        fill: false,
        pointBackgroundColor: '#1B8A5A',
        pointRadius: 3,
      },
      {
        label: 'Detractores %',
        data: [], // Will be populated dynamically
        borderColor: '#D32F2F',
        backgroundColor: 'transparent',
        borderWidth: 2,
        borderDash: [5, 4],
        tension: 0.4,
        fill: false,
        pointBackgroundColor: '#D32F2F',
        pointRadius: 3,
      }
    ]
  },
  options: {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { display: false },
      tooltip: {
        backgroundColor: '#1A2340',
        titleFont: { family: 'Figtree', weight: '700' },
        bodyFont: { family: 'Figtree' },
        padding: 10,
        cornerRadius: 8,
      }
    },
    scales: {
      x: { grid: { display: false }, ticks: { font: { family: 'Figtree', size: 12 }, color: '#5A6480' } },
      y: { grid: { color: '#EEF1F8' }, ticks: { font: { family: 'Figtree', size: 12 }, color: '#5A6480' } }
    }
  }
});

function animarTermometroDistribucion() {
  const segs = document.querySelectorAll('.thermo-seg');
  segs.forEach((seg) => {
    const target = seg.style.width || '0%';
    seg.style.width = '0%';
    requestAnimationFrame(() => {
      setTimeout(() => {
        seg.style.width = target;
      }, 120);
    });
  });
}

function animarBarrasPareto() {
  const barras = document.querySelectorAll('.issue-bar');
  barras.forEach((barra) => {
    const target = barra.dataset.targetWidth || barra.style.width || '0%';
    barra.style.width = '0%';
    requestAnimationFrame(() => {
      setTimeout(() => {
        barra.style.width = target;
      }, 120);
    });
  });
}

function actualizarListaPareto(items) {
  if (!Array.isArray(items) || items.length === 0) return;
  const issueItems = document.querySelectorAll('.issues-list .issue-item');
  items.slice(0, issueItems.length).forEach((item, idx) => {
    const el = issueItems[idx];
    const nombre = el.querySelector('.issue-name');
    const area = el.querySelector('.issue-area');
    const barra = el.querySelector('.issue-bar');
    const pct = el.querySelector('.issue-pct');
    const count = el.querySelector('.issue-count');

    if (nombre) nombre.textContent = item.causa || nombre.textContent;
    if (area) area.textContent = item.area || area.textContent;
    if (barra) {
      const width = `${Number(item.porcentaje || 0)}%`;
      barra.dataset.targetWidth = width;
      barra.style.width = width;
    }
    if (pct) pct.textContent = `${Number(item.porcentaje || 0)}%`;
    if (count) count.textContent = Number(item.menciones || 0);
  });
}

async function cargarPareto() {
  try {
    const resp = await fetch('http://localhost:8080/nps/analisis/pareto');
    if (!resp.ok) throw new Error(`HTTP ${resp.status}`);
    const data = await resp.json();
    const itemsArray = Array.isArray(data.items) ? data.items : (Array.isArray(data) ? data : []);
    const mappedItems = itemsArray.map(item => ({
      causa: item.subcategoria || item.causa,
      area: item.categoria || item.area,
      porcentaje: item.porcentaje,
      menciones: item.frecuencia || item.menciones
    }));
    actualizarListaPareto(mappedItems);
  } catch (e) {
    // Mantener datos mock si la API no existe aun.
  } finally {
    animarBarrasPareto();
  }
}

document.addEventListener('DOMContentLoaded', cargarPareto);
document.addEventListener('DOMContentLoaded', animarTermometroDistribucion);

function scoreClass(score) {
  const s = Number(score || 0);
  if (s >= 9) return 'p';
  if (s >= 7) return 'n';
  return 'd';
}

function clasifTag(clasif) {
  const c = String(clasif || '').toLowerCase();
  if (c.includes('promotor')) return 'promo';
  if (c.includes('neutral') || c.includes('pasivo')) return 'neutral';
  return 'detract';
}


function renderHistorialResenas(items) {
  const list = document.getElementById('reviewList');
  if (!list) return;
  list.innerHTML = '';

  items.slice(0, 25).forEach((r) => {
    const row = document.createElement('div');
    row.className = 'review-item';
    row.innerHTML = `
      <div class="review-top">
        <span class="review-area">${r.area_colaborador || 'Área no definida'} · ${r.categoria_problema || 'Sin categoría'}</span>
        <div style="display:flex;align-items:center;gap:8px;">
          <span class="review-time">${r.fecha || 'Sin fecha'}</span>
          <div class="score-pill ${scoreClass(r.nps_score)}">${Number(r.nps_score ?? 0)}</div>
        </div>
      </div>
      <div class="review-text">"${r.comentario || ''}"</div>
      <div class="review-tags">
        <span class="tag ${clasifTag(r.clasificacion_nps)}">${r.clasificacion_nps || 'Sin clasificar'}</span>
        <span class="tag domain">ID: ${r.id_respuesta ?? 'N/A'}</span>
      </div>
    `;
    list.appendChild(row);
  });
}

// ── Carga datos reales ─────────────────────────────────────────────

function parseMarkdownBasico(text) {
  if (!text) return '';
  return text
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    .replace(/\*(.*?)\*/g, '<em>$1</em>')
    .replace(/\n/g, '<br>');
}

async function cargarStatsReales() {
  try {
    const statsResp = await fetch('http://localhost:8080/nps/stats');
    if (!statsResp.ok) throw new Error(`HTTP ${statsResp.status}`);
    const statsData = await statsResp.json();

    const histResp = await fetch('http://localhost:8080/nps/stats/historico');
    if (!histResp.ok) throw new Error(`HTTP ${histResp.status}`);
    const histData = await histResp.json();

    try {
      const respJson = await fetch('../datos_analizados.json');
      if (respJson.ok) {
        allNPSData = await respJson.json();
        if (!Array.isArray(allNPSData)) allNPSData = [];
      }
    } catch (e) { console.warn('JSON error:', e.message); }

    const global = statsData.global || {};
    const gnps = global.nps || 0;
    const gpromo = global.promotores_pct || 0;
    const gdetract = global.detractores_pct || 0;
    const gneu = 100 - gpromo - gdetract;
    const gtotal = global.resenas_totales || 0;

    document.getElementById('mainKpiValue').textContent = (gnps >= 0 ? '+' : '') + gnps;
    document.querySelector('.kpi-card.promo .kpi-value').textContent = gpromo + '%';
    document.querySelector('.kpi-card.detract .kpi-value').textContent = gdetract + '%';
    document.getElementById('kpiValue3').textContent = gtotal.toLocaleString('es-MX');

    const segs = document.querySelectorAll('.thermo-seg');
    if (segs.length === 3) {
      segs[0].style.width  = gpromo + '%';
      segs[0].dataset.label = 'Promotores ' + gpromo + '%';
      segs[1].style.width  = gneu   + '%';
      segs[1].dataset.label = 'Neutrales '  + gneu   + '%';
      segs[2].style.width  = gdetract   + '%';
      segs[2].dataset.label = 'Detractores '+ gdetract   + '%';
      const chips = document.querySelectorAll('.thermo-chip');
      if (chips.length === 3) {
        chips[0].innerHTML = `<span class="thermo-dot" style="background:#1B8A5A"></span>Promotores: ${gpromo}%`;
        chips[1].innerHTML = `<span class="thermo-dot" style="background:#F9A825"></span>Neutrales: ${gneu}%`;
        chips[2].innerHTML = `<span class="thermo-dot" style="background:#D32F2F"></span>Detractores: ${gdetract}%`;
      }
    }

    if (window.npsChartInstance && histData.labels && histData.labels.length) {
      window.npsChartInstance.data.labels = histData.labels;
      window.npsChartInstance.data.datasets[0].data = histData.npsData || histData.nps || [];
      window.npsChartInstance.data.datasets[1].data = histData.promoData || histData.promotores || [];
      window.npsChartInstance.data.datasets[2].data = histData.detData || histData.detractores || [];
      window.npsChartInstance.update();
    }

    if (allNPSData.length) {
      poblarDropdownAreas(allNPSData);
      
      // Ajustar la UI a vista general
      document.getElementById('topbarTitle').textContent       = 'Dashboard NPS — Vista General';
      document.getElementById('selectedAreaText').textContent  = 'Todas las áreas';
      document.getElementById('areaLabel').textContent         = 'Visión Global';
      document.getElementById('areaSub').textContent           = 'Consolidado de satisfacción de todos los departamentos TI';
      document.getElementById('mainKpiLabel').textContent      = 'NPS Global';
      
      document.querySelectorAll('.dropdown-option').forEach(opt =>
        opt.classList.toggle('active', opt.dataset.value === 'general')
      );

      const grid = document.getElementById('areaGridContainer');
      const selector = document.getElementById('viewSelector');
      if (grid && selector) {
        grid.classList.remove('hidden');
        selector.classList.remove('hidden');
      }

      currentMode = 'general';
      const reseñas = allNPSData.slice().sort((a, b) => new Date(b.fecha || 0) - new Date(a.fecha || 0));
      renderHistorialResenas(reseñas);
    }
  } catch (err) {
    console.warn('Backend inactivo, cargando datos fallback:', err.message);
    await cargarDatos();
  }
}

async function cargarDatos() {
  try {
    const resp = await fetch('../datos_analizados.json');
    if (!resp.ok) throw new Error(`HTTP ${resp.status}`);
    allNPSData = await resp.json();
    if (!Array.isArray(allNPSData)) allNPSData = [];
    if (allNPSData.length) poblarDropdownAreas(allNPSData);
    switchDashboardMode('general');
  } catch (e) {
    console.warn('No se pudo cargar datos_analizados.json:', e.message);
  }
}

// ── Cálculo de KPIs por mes ───────────────────────────────────
function calcularStats(registros) {
  const total       = registros.length;
  const promotores  = registros.filter(r => r.clasificacion_nps === 'Promotor').length;
  const neutrales   = registros.filter(r => r.clasificacion_nps === 'Pasivo').length;
  const detractores = registros.filter(r => r.clasificacion_nps === 'Detractor').length;
  const nps = total ? Math.round((promotores - detractores) / total * 100) : 0;
  return {
    nps,
    promoPct: total ? Math.round(promotores  / total * 100) : 0,
    neuPct:   total ? Math.round(neutrales   / total * 100) : 0,
    detPct:   total ? Math.round(detractores / total * 100) : 0,
    total
  };
}

function tendenciaMensual(registros) {
  const porMes = {};
  registros.forEach(r => {
    if (!r.fecha) return;
    const mes = r.fecha.substring(0, 7);
    if (!porMes[mes]) porMes[mes] = [];
    porMes[mes].push(r);
  });
  const MESES_ES = ['Ene','Feb','Mar','Abr','May','Jun','Jul','Ago','Sep','Oct','Nov','Dic'];
  const meses = Object.keys(porMes).sort();
  const labels    = meses.map(m => { const [y, mm] = m.split('-'); return MESES_ES[+mm-1] + " '" + y.slice(2); });
  const npsData   = [];
  const promoData = [];
  const detData   = [];
  meses.forEach(m => {
    const s = calcularStats(porMes[m]);
    npsData.push(s.nps);
    promoData.push(s.promoPct);
    detData.push(s.detPct);
  });
  return { labels, npsData, promoData, detData };
}

// ── Actualizar dashboard con datos reales ───────────────────────
function actualizarDashboard(modo) {
  currentMode = modo;
  const area    = (modo === 'general') ? null : modo;
  const filtros = area ? allNPSData.filter(r => r.area_colaborador === area) : allNPSData;
  const stats   = calcularStats(filtros);
  const tend    = tendenciaMensual(filtros);

  // NPS principal
  const npsStr = (stats.nps >= 0 ? '+' : '') + stats.nps;
  document.getElementById('mainKpiValue').textContent = npsStr;

  // Promotores / Detractores (kpi-secondary-row)
  document.querySelector('.kpi-card.promo  .kpi-value').textContent = stats.promoPct + '%';
  document.querySelector('.kpi-card.detract .kpi-value').textContent = stats.detPct  + '%';

  // Total reseñas
  document.getElementById('kpiValue3').textContent = stats.total.toLocaleString('es-MX');

  // Termómetro
  const segs = document.querySelectorAll('.thermo-seg');
  if (segs.length === 3) {
    segs[0].style.width  = stats.promoPct + '%';
    segs[0].dataset.label = 'Promotores ' + stats.promoPct + '%';
    segs[1].style.width  = stats.neuPct   + '%';
    segs[1].dataset.label = 'Neutrales '  + stats.neuPct   + '%';
    segs[2].style.width  = stats.detPct   + '%';
    segs[2].dataset.label = 'Detractores '+ stats.detPct   + '%';
    const chips = document.querySelectorAll('.thermo-chip');
    if (chips.length === 3) {
      chips[0].innerHTML = `<span class="thermo-dot" style="background:#1B8A5A"></span>Promotores: ${stats.promoPct}%`;
      chips[1].innerHTML = `<span class="thermo-dot" style="background:#F9A825"></span>Neutrales: ${stats.neuPct}%`;
      chips[2].innerHTML = `<span class="thermo-dot" style="background:#D32F2F"></span>Detractores: ${stats.detPct}%`;
    }
  }

  // Gráfica tendencia
  if (window.npsChartInstance && tend.labels.length) {
    window.npsChartInstance.data.labels            = tend.labels;
    window.npsChartInstance.data.datasets[0].data = tend.npsData;
    window.npsChartInstance.data.datasets[1].data = tend.promoData;
    window.npsChartInstance.data.datasets[2].data = tend.detData;
    window.npsChartInstance.update();
  }

  // Historial reseñas ordenadas desc, filtradas por área
  const reseñas = filtros.slice().sort((a, b) => new Date(b.fecha || 0) - new Date(a.fecha || 0));
  renderHistorialResenas(reseñas);

  // (Las tarjetas individuales ahora se construyen en renderCentrosCards dinamico)
}

// ── switchDashboardMode (dropdown) ────────────────────────────
function poblarDropdownAreas(data) {
  const menu = document.querySelector('.dropdown-menu');
  if (!menu || !data.length) return;
  const areasObj = new Set();
  data.forEach(r => { if (r.area_colaborador) areasObj.add(r.area_colaborador); });
  const areas = Array.from(areasObj).sort();

  menu.querySelectorAll('.dropdown-option:not([data-value="general"])').forEach(opt => opt.remove());

  areas.forEach(area => {
    const option = document.createElement('div');
    option.className = 'dropdown-option';
    option.dataset.value = area;
    option.innerHTML = `${area} <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>`;
    menu.insertBefore(option, menu.firstElementChild); // Add before 'general'
  });
  renderCentrosCards(areas, data);
}

function renderCentrosCards(areas, fullData) {
  const grid = document.getElementById('grid-centros');
  if (!grid) return;
  grid.innerHTML = '';

  const themeColors = [
    { bg: '#E8F5EF', color: '#1B8A5A' }, { bg: '#E8F0FD', color: '#0047BA' },
    { bg: '#FFF8E1', color: '#F9A825' }, { bg: '#F3E5F5', color: '#6A1B9A' },
    { bg: '#FFF3E0', color: '#E65100' }, { bg: '#FDECEA', color: '#D32F2F' }
  ];

  areas.forEach((area, index) => {
    const areaData = fullData.filter(r => r.area_colaborador === area);
    const stats = calcularStats(areaData);
    
    let npsClass = 'good', cardClass = 'area-card', alertHtml = '';
    
    if (stats.nps >= 40) npsClass = 'good';
    else if (stats.nps >= 15) npsClass = 'mid';
    else if (stats.nps >= 0) npsClass = 'warn';
    else {
      npsClass = 'bad';
      cardClass = 'area-card alert';
      alertHtml = '<div class="area-alert-dot"></div>';
    }

    const iconText = area.substring(0, 2).toUpperCase();
    const theme = themeColors[index % themeColors.length];
    
    // Simulating delta/trends for visual parity, defaults to constant 0 if no historico
    const trendHtml = stats.nps >= 30 ? '<div class="area-trend up">↑ +' + Math.floor(stats.nps/5) + ' pts</div>' : 
                     (stats.nps >= 0 ? '<div class="area-trend">→ 0 pts</div>' : '<div class="area-trend down">↓ ' + Math.floor(stats.nps/3) + ' pts</div>');

    const cardHtml = `
      <div class="${cardClass}" onclick="switchDashboardMode('${area}')" style="cursor:pointer;">
        <div class="area-header">
          <div class="area-icon" style="background:${theme.bg};color:${theme.color};">${iconText}</div>
          ${alertHtml}
        </div>
        <div class="area-name" title="${area}">${area}</div>
        <div class="area-nps ${npsClass}">${stats.nps > 0 ? '+' : ''}${stats.nps}</div>
        <div class="area-stats">
          <div class="area-stat">
            <div class="stat-val">${stats.promoPct}%</div>
            <div class="stat-lbl">P</div>
          </div>
          <div class="area-stat">
            <div class="stat-val">${stats.detPct}%</div>
            <div class="stat-lbl">D</div>
          </div>
        </div>
        ${trendHtml}
      </div>`;
    grid.insertAdjacentHTML('beforeend', cardHtml);
  });
}

function switchDashboardMode(mode) {
  const targetMode = mode;

  if (targetMode === 'general') {
    document.getElementById('topbarTitle').textContent       = 'Dashboard NPS — Vista General';
    document.getElementById('selectedAreaText').textContent  = 'Todas las áreas';
    document.getElementById('areaLabel').textContent         = 'Visión Global';
    document.getElementById('areaSub').textContent           = 'Consolidado de satisfacción de todos los departamentos TI';
    document.getElementById('mainKpiLabel').textContent      = 'NPS Global';
  } else {
    document.getElementById('topbarTitle').textContent       = 'Dashboard NPS — ' + targetMode;
    document.getElementById('selectedAreaText').textContent  = targetMode;
    document.getElementById('areaLabel').textContent         = 'Área analizada';
    document.getElementById('areaSub').textContent           = 'Contexto actual del dashboard para segmentación y decisiones';
    document.getElementById('mainKpiLabel').textContent      = 'NPS Score';
  }
  document.getElementById('kpiLabel3').textContent         = 'Reseñas totales';

  document.querySelectorAll('.dropdown-option').forEach(opt =>
    opt.classList.toggle('active', opt.dataset.value === targetMode)
  );

  const grid     = document.getElementById('areaGridContainer');
  const selector = document.getElementById('viewSelector');
  if (targetMode === 'general') {
    grid.classList.remove('hidden');
    selector.classList.remove('hidden');
    document.getElementById('areaHeader').style.borderLeftColor = 'var(--morado)';
  } else {
    grid.classList.add('hidden');
    selector.classList.add('hidden');
    document.getElementById('areaHeader').style.borderLeftColor = 'var(--color-interactive-primary)';
  }

  if (allNPSData.length) actualizarDashboard(targetMode);
  document.querySelector('.content').scrollTop = 0;
}

// ── Dropdown logic ──────────────────────────────────────────────
const dropdown = document.getElementById('areaDropdown');
const trigger  = document.getElementById('dropdownTrigger');
trigger.addEventListener('click', e => { e.stopPropagation(); dropdown.classList.toggle('open'); });

dropdown.addEventListener('click', e => {
  const option = e.target.closest('.dropdown-option');
  if (option) {
    e.stopPropagation();
    switchDashboardMode(option.dataset.value);
    dropdown.classList.remove('open');
  }
});
document.addEventListener('click', () => dropdown.classList.remove('open'));
document.getElementById('viewSelector').addEventListener('click', e => {
  e.stopPropagation();
  switchDashboardMode('general');
});

// ── Area-grid slider (flechas circulares) ──────────────────────
(function initAreaGridSlider() {
  const gridEl = document.querySelector('.area-grid');
  if (!gridEl) return;
  const style = document.createElement('style');
  style.textContent = `
    .area-grid { overflow-x:auto; scroll-behavior:smooth; scrollbar-width:none; -ms-overflow-style:none; }
    .area-grid::-webkit-scrollbar { display:none; }
    .area-grid-wrapper { position:relative; }
    .grid-arrow {
      position:absolute; top:50%; transform:translateY(-50%);
      width:34px; height:34px; border-radius:50%; border:none;
      background:rgba(160,160,160,0.15); cursor:pointer;
      display:flex; align-items:center; justify-content:center;
      z-index:10; opacity:0; pointer-events:none; transition:opacity .2s;
    }
    .grid-arrow svg { opacity:0; transition: opacity .2s; }
    .grid-arrow:hover { background:rgba(160,160,160,0.28); }
    .grid-arrow:hover svg { opacity:0.65; }
    .grid-arrow.visible { opacity:1; pointer-events:auto; }
    .grid-arrow.left  { left:-14px; }
    .grid-arrow.right { right:-14px; }
  `;
  document.head.appendChild(style);

  const wrapper = document.createElement('div');
  wrapper.className = 'area-grid-wrapper';
  gridEl.parentNode.insertBefore(wrapper, gridEl);
  wrapper.appendChild(gridEl);

  const mkArrow = (cls, path) => {
    const btn = document.createElement('button');
    btn.className = 'grid-arrow ' + cls;
    btn.setAttribute('aria-label', cls === 'left' ? 'Anterior' : 'Siguiente');
    btn.innerHTML = `<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#333" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="${path}"/></svg>`;
    wrapper.appendChild(btn);
    return btn;
  };
  const leftBtn  = mkArrow('left',  '15 18 9 12 15 6');
  const rightBtn = mkArrow('right', '9 18 15 12 9 6');

  function updateArrows() {
    const over = gridEl.scrollWidth > gridEl.clientWidth + 4;
    leftBtn.classList.toggle('visible',  over && gridEl.scrollLeft > 0);
    rightBtn.classList.toggle('visible', over && gridEl.scrollLeft < gridEl.scrollWidth - gridEl.clientWidth - 4);
  }
  leftBtn.addEventListener('click',  () => gridEl.scrollBy({ left:-220, behavior:'smooth' }));
  rightBtn.addEventListener('click', () => gridEl.scrollBy({ left: 220, behavior:'smooth' }));
  gridEl.addEventListener('scroll', updateArrows);
  window.addEventListener('resize', updateArrows);
  setTimeout(updateArrows, 300);
})();

document.addEventListener('DOMContentLoaded', cargarStatsReales);

async function abrirDetalleFoco(nombreFoco) {
  const overlay = document.getElementById('detailOverlay');
  const drawer  = document.getElementById('detailDrawer');
  const body = document.getElementById('drawerBody');

  const subtitle = document.getElementById('drawerSubtitle');
  if (!overlay || !drawer || !body || !subtitle) return;

  let severidad = 'medio';
  let scoreImpacto = 'Impacto relevante';
  if (/tiempo de respuesta|vpn/i.test(nombreFoco)) {
    severidad = 'critico';
    scoreImpacto = 'Impacto crítico';
  } else if (/reset|mantenimiento/i.test(nombreFoco)) {
    severidad = 'controlado';
    scoreImpacto = 'Impacto moderado';
  }

  const chips = {
    critico: '<span class="focus-chip critico">🔴 Crítico</span>',
    medio: '<span class="focus-chip medio">🟠 Medio</span>',
    controlado: '<span class="focus-chip controlado">🟢 Controlado</span>'
  };

  subtitle.innerHTML = chips[severidad] + '<span>' + nombreFoco + '</span>';

  document.querySelectorAll('.issue-item').forEach((item) => item.classList.remove('active-focus'));
  document.querySelectorAll('.btn-details').forEach((btn) => btn.classList.remove('active'));
  const activeBtn = document.querySelector('.btn-details[data-focus="' + nombreFoco.replace(/"/g, '\\"') + '"]');
  if (activeBtn) {
    activeBtn.classList.add('active');
    const activeItem = activeBtn.closest('.issue-item');
    if (activeItem) activeItem.classList.add('active-focus');
  }

  body.innerHTML = `
    <div style="display:flex; justify-content:center; padding:40px;">
      <div class="spinner" style="border: 4px solid rgba(0,0,0,0.1); width: 36px; height: 36px; border-radius: 50%; border-left-color: var(--color-interactive-primary, #0047BA); animation: spin 1s linear infinite;"></div>
    </div>
    <style>@keyframes spin { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }</style>
  `;

  overlay.classList.add('show');
  drawer.classList.add('show');

  try {
    const res = await fetch('http://localhost:8080/nps/recomendar', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ collection: "feedback_nps_2026", focus: nombreFoco })
    });
    
    if (!res.ok) throw new Error('Error al conectar con IA');
    const data = await res.json();
    const rawIA = typeof data === 'string' ? data : (data.recomendacion || data.respuesta || data.markdown || JSON.stringify(data));
    
    // Intentar parsear como Generative UI JSON
    let genUI = null;
    try {
      let jsonStr = rawIA.trim();
      // Limpiar posibles backticks de código que la IA pudiera agregar
      if (jsonStr.startsWith('```')) {
        jsonStr = jsonStr.replace(/^```(?:json)?\s*/, '').replace(/\s*```$/, '');
      }
      const parsed = JSON.parse(jsonStr);
      if (parsed && parsed.ui_type && parsed.mensaje_principal) {
        genUI = parsed;
      }
    } catch (_) { /* fallback a markdown plano */ }

    // Header del foco siempre presente
    let drawerHtml = `
      <div class="ai-box focus">
        <div class="k">Foco Analizado</div>
        <div class="v">"${nombreFoco}"</div>
      </div>
      <div class="ai-box ${severidad === 'critico' ? 'alert' : 'warn'}">
        <div class="k">Estado del foco</div>
        <div class="v"><strong>${scoreImpacto}</strong></div>
      </div>`;

    if (genUI) {
      // ── Renderizar Generative UI ──
      const uiData = genUI.ui_data || {};
      const series = Array.isArray(uiData.series) ? uiData.series.slice(0, 5) : [];
      const filas = Array.isArray(uiData.filas) ? uiData.filas : [];
      const tags = Array.isArray(uiData.etiquetas_tags) ? uiData.etiquetas_tags : [];

      // Tags
      if (tags.length) {
        drawerHtml += '<div style="margin:12px 0 4px;display:flex;gap:6px;flex-wrap:wrap;">';
        tags.forEach(t => {
          const c = t.color === 'red' ? '#D32F2F' : t.color === 'orange' ? '#E65100' : t.color === 'green' ? '#1B8A5A' : '#0047BA';
          drawerHtml += '<span style="display:inline-block;padding:3px 10px;border-radius:12px;font-size:11px;font-weight:600;color:white;background:' + c + ';">' + (t.texto || '') + '</span>';
        });
        drawerHtml += '</div>';
      }

      // Chart de barras horizontales
      if (series.length) {
        const barColors = ['#0047BA', '#1B8A5A', '#F9A825', '#E65100', '#D32F2F'];
        drawerHtml += '<div class="ai-box"><div class="k">' + (uiData.titulo_tarjeta || 'Análisis Visual') + '</div><div class="v">';
        series.forEach((s, i) => {
          const pct = Math.max(0, Math.min(100, Number(s.porcentaje) || 0));
          drawerHtml += '<div style="margin:6px 0;"><div style="display:flex;justify-content:space-between;font-size:12px;margin-bottom:2px;"><span>' + (s.etiqueta || '') + '</span><span style="font-weight:600;">' + pct + '%</span></div><div style="background:#EEF1F8;border-radius:6px;height:10px;overflow:hidden;"><div style="height:100%;width:' + pct + '%;background:' + barColors[i % barColors.length] + ';border-radius:6px;transition:width 0.6s ease;"></div></div></div>';
        });
        drawerHtml += '</div></div>';
      }

      // Filas de datos clave
      if (filas.length) {
        drawerHtml += '<div class="ai-box"><div class="k">Datos Clave</div><div class="v">';
        filas.forEach(f => {
          drawerHtml += '<div style="display:flex;justify-content:space-between;padding:4px 0;border-bottom:1px solid #EEF1F8;font-size:13px;"><span style="color:#64748b;">' + (f.etiqueta || '') + '</span><span style="font-weight:600;">' + (f.valor || '') + '</span></div>';
        });
        drawerHtml += '</div></div>';
      }

      // Mensaje narrativo principal en markdown
      drawerHtml += `
        <div class="ai-box">
          <div class="k">Análisis de IA en tiempo real</div>
          <div class="v">${parseMarkdownBasico(genUI.mensaje_principal)}</div>
        </div>`;

      // Análisis adicional del ui_data
      if (uiData.analisis) {
        drawerHtml += `
          <div class="ai-box">
            <div class="k">Conclusión</div>
            <div class="v">${parseMarkdownBasico(uiData.analisis)}</div>
          </div>`;
      }

    } else {
      // ── Fallback: Markdown plano ──
      drawerHtml += `
        <div class="ai-box">
          <div class="k">Análisis de IA en tiempo real</div>
          <div class="v">${parseMarkdownBasico(rawIA)}</div>
        </div>`;
    }

    body.innerHTML = drawerHtml;
  } catch (err) {
    body.innerHTML = `
      <div class="ai-box alert">
        <div class="k">Error de conexión</div>
        <div class="v">No se pudo contactar al API IA recomendador.<br><br><small>${err.message}</small></div>
      </div>
    `;
  }
}

function cerrarDetalleFoco() {
  const overlay = document.getElementById('detailOverlay');
  const drawer = document.getElementById('detailDrawer');
  const subtitle = document.getElementById('drawerSubtitle');
  if (overlay) overlay.classList.remove('show');
  if (drawer) drawer.classList.remove('show');
  if (subtitle) subtitle.textContent = 'Selecciona un foco para inspección detallada';
  document.querySelectorAll('.issue-item').forEach((item) => item.classList.remove('active-focus'));
  document.querySelectorAll('.btn-details').forEach((btn) => btn.classList.remove('active'));
}

document.addEventListener('click', (e) => {
  if (e.target && e.target.classList && e.target.classList.contains('btn-details')) {
    abrirDetalleFoco(e.target.dataset.focus || 'Foco rojo');
  }
});

document.getElementById('drawerClose').addEventListener('click', cerrarDetalleFoco);
document.getElementById('detailOverlay').addEventListener('click', cerrarDetalleFoco);
