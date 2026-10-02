// Card Reveal Challenge - Client Script & Sound Effects

// Inisialisasi Audio FX menggunakan Web Audio API murni (tanpa file audio eksternal)
class SoundFX {
  constructor() {
    this.ctx = null;
  }

  init() {
    if (!this.ctx) {
      const AudioCtx = window.AudioContext || window.webkitAudioContext;
      if (AudioCtx) {
        this.ctx = new AudioCtx();
      }
    }
    if (this.ctx && this.ctx.state === 'suspended') {
      this.ctx.resume();
    }
  }

  playFlip() {
    this.init();
    if (!this.ctx) return;
    try {
      const osc = this.ctx.createOscillator();
      const gain = this.ctx.createGain();
      osc.type = 'sine';
      osc.frequency.setValueAtTime(280, this.ctx.currentTime);
      osc.frequency.exponentialRampToValueAtTime(700, this.ctx.currentTime + 0.18);
      gain.gain.setValueAtTime(0.2, this.ctx.currentTime);
      gain.gain.exponentialRampToValueAtTime(0.01, this.ctx.currentTime + 0.18);
      osc.connect(gain);
      gain.connect(this.ctx.destination);
      osc.start();
      osc.stop(this.ctx.currentTime + 0.18);
    } catch (e) {
      console.warn('Audio error:', e);
    }
  }

  playFanfare() {
    this.init();
    if (!this.ctx) return;
    try {
      const notes = [523.25, 659.25, 783.99, 1046.50]; // Chord C - E - G - C tinggi
      notes.forEach((freq, idx) => {
        const osc = this.ctx.createOscillator();
        const gain = this.ctx.createGain();
        osc.type = 'triangle';
        osc.frequency.setValueAtTime(freq, this.ctx.currentTime + idx * 0.08);
        gain.gain.setValueAtTime(0.25, this.ctx.currentTime + idx * 0.08);
        gain.gain.exponentialRampToValueAtTime(0.001, this.ctx.currentTime + idx * 0.08 + 0.35);
        osc.connect(gain);
        gain.connect(this.ctx.destination);
        osc.start(this.ctx.currentTime + idx * 0.08);
        osc.stop(this.ctx.currentTime + idx * 0.08 + 0.35);
      });
    } catch (e) {
      console.warn('Fanfare error:', e);
    }
  }
}

const sfx = new SoundFX();

// Partikel Confetti Sederhana Canvas
function launchConfetti() {
  const canvas = document.getElementById('confettiCanvas');
  if (!canvas) return;
  const ctx = canvas.getContext('2d');
  canvas.width = window.innerWidth;
  canvas.height = window.innerHeight;

  const pieces = [];
  const colors = ['#f59e0b', '#ef4444', '#3b82f6', '#10b981', '#8b5cf6', '#ec4899'];
  for (let i = 0; i < 70; i++) {
    pieces.push({
      x: canvas.width / 2,
      y: canvas.height / 2,
      vx: (Math.random() - 0.5) * 14,
      vy: (Math.random() - 0.8) * 16,
      size: Math.random() * 8 + 5,
      color: colors[Math.floor(Math.random() * colors.length)],
      alpha: 1,
      rotation: Math.random() * 360,
      vRot: (Math.random() - 0.5) * 10
    });
  }

  let frame = 0;
  function render() {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    let active = false;
    for (const p of pieces) {
      p.x += p.vx;
      p.y += p.vy;
      p.vy += 0.35; // gravitasi
      p.vx *= 0.98;
      p.rotation += p.vRot;
      p.alpha -= 0.012;

      if (p.alpha > 0) {
        active = true;
        ctx.save();
        ctx.globalAlpha = p.alpha;
        ctx.translate(p.x, p.y);
        ctx.rotate((p.rotation * Math.PI) / 180);
        ctx.fillStyle = p.color;
        ctx.fillRect(-p.size / 2, -p.size / 2, p.size, p.size);
        ctx.restore();
      }
    }

    if (active && frame < 120) {
      frame++;
      requestAnimationFrame(render);
    } else {
      ctx.clearRect(0, 0, canvas.width, canvas.height);
    }
  }
  requestAnimationFrame(render);
}

// Logika Klik Kartu
document.addEventListener('DOMContentLoaded', () => {
  const cardElements = document.querySelectorAll('.card-item');
  const modal = document.getElementById('challengeModal');
  const modalNumber = document.getElementById('modalCardNumber');
  const modalSuit = document.getElementById('modalCardSuit');
  const modalText = document.getElementById('modalChallengeText');
  const closeModalBtn = document.getElementById('closeModalBtn');
  const resetBtn = document.getElementById('resetSessionBtn');

  // Event Klik pada setiap kartu
  cardElements.forEach(card => {
    card.addEventListener('click', async () => {
      const cardId = card.dataset.id;
      const isRevealed = card.dataset.revealed === 'true';
      const inner = card.querySelector('.card-flip-inner');
      const challengeText = card.dataset.challenge;
      const cardNumber = card.dataset.number;
      const suitSymbol = card.dataset.symbol;
      const suitColor = card.dataset.color;

      // Jika sudah terbuka, buka langsung modil fokusnya
      if (isRevealed) {
        showChallengeModal(cardNumber, suitSymbol, suitColor, challengeText);
        return;
      }

      // Suara Flip
      sfx.playFlip();

      // Putar animasi kartu di grid
      inner.classList.add('is-flipped');
      card.dataset.revealed = 'true';
      card.classList.add('is-revealed-badge');

      // Tampilkan tanda centang / badge selesai
      const badge = card.querySelector('.revealed-tag');
      if (badge) {
        badge.classList.remove('hidden');
      }

      // Kirim status ke server
      try {
        await fetch(`/api/cards/${cardId}/reveal`, {
          method: 'POST'
        });
      } catch (err) {
        console.error('Gagal sinkron status reveal:', err);
      }

      // Jeda 500ms agar animasi flip selesai sebelum modal muncul
      setTimeout(() => {
        sfx.playFanfare();
        launchConfetti();
        showChallengeModal(cardNumber, suitSymbol, suitColor, challengeText);
      }, 550);
    });
  });

  // Tampilkan Modal Tantangan
  function showChallengeModal(number, symbol, color, text) {
    if (!modal) return;
    modalNumber.textContent = `KARTU #${number}`;
    modalSuit.textContent = symbol;
    modalSuit.className = `text-5xl font-black ${color}`;
    modalText.textContent = text;
    modal.classList.remove('hidden');
    modal.classList.add('flex');
  }

  // Tutup Modal
  function hideModal() {
    if (!modal) return;
    modal.classList.add('hidden');
    modal.classList.remove('flex');
  }

  if (closeModalBtn) {
    closeModalBtn.addEventListener('click', hideModal);
  }

  if (modal) {
    modal.addEventListener('click', (e) => {
      if (e.target === modal) {
        hideModal();
      }
    });
  }

  // Keyboard shortcut Esc untuk tutup modal
  window.addEventListener('keydown', (e) => {
    if (e.key === 'Escape' && modal && !modal.classList.contains('hidden')) {
      hideModal();
    }
  });

  // Tombol Reset Sesi
  if (resetBtn) {
    resetBtn.addEventListener('click', async () => {
      const confirmReset = confirm('Mulai sesi baru? Semua kartu akan ditutup kembali.');
      if (!confirmReset) return;

      try {
        const res = await fetch('/api/cards/reset', {
          method: 'POST'
        });
        const data = await res.json();
        if (data.success) {
          // Putar balik semua kartu
          cardElements.forEach((card, idx) => {
            const inner = card.querySelector('.card-flip-inner');
            setTimeout(() => {
              inner.classList.remove('is-flipped');
              card.dataset.revealed = 'false';
              card.classList.remove('is-revealed-badge');
              const badge = card.querySelector('.revealed-tag');
              if (badge) {
                badge.classList.add('hidden');
              }
            }, idx * 30); // Stagger animation
          });
        }
      } catch (err) {
        alert('Gagal mereset kartu: ' + err.message);
      }
    });
  }
});

