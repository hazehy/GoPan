<template>
  <canvas ref="canvasRef" class="polygon-bg-canvas" />
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from "vue";

interface PolygonParticle {
  x: number;
  y: number;
  radius: number;
  sides: number;
  angle: number;
  rotateSpeed: number;
  vx: number;
  vy: number;
  color: string;
  alpha: number;
}

const canvasRef = ref<HTMLCanvasElement | null>(null);

let particles: PolygonParticle[] = [];
let frameId = 0;
let width = 0;
let height = 0;
let dpr = 1;

const mouse = {
  x: -9999,
  y: -9999,
  active: false,
};

const COLORS = ["#ffffff", "#dbeafe", "#bfdbfe", "#93c5fd", "#60a5fa"];
const PARTICLE_COUNT = 84;
const ATTRACT_RADIUS = 160;

function rand(min: number, max: number) {
  return Math.random() * (max - min) + min;
}

function createParticle(): PolygonParticle {
  return {
    x: rand(0, width),
    y: rand(0, height),
    radius: rand(10, 26),
    sides: Math.floor(rand(3, 7)),
    angle: rand(0, Math.PI * 2),
    rotateSpeed: rand(-0.0035, 0.0035),
    vx: rand(-0.35, 0.35),
    vy: rand(-0.35, 0.35),
    color: COLORS[Math.floor(rand(0, COLORS.length))],
    alpha: rand(0.22, 0.55),
  };
}

function resizeCanvas() {
  const canvas = canvasRef.value;
  if (!canvas) {
    return;
  }

  const rect = canvas.getBoundingClientRect();
  width = rect.width;
  height = rect.height;
  dpr = Math.max(window.devicePixelRatio || 1, 1);

  canvas.width = Math.floor(width * dpr);
  canvas.height = Math.floor(height * dpr);

  const ctx = canvas.getContext("2d");
  if (!ctx) {
    return;
  }
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0);

  particles = Array.from({ length: PARTICLE_COUNT }, () => createParticle());
}

function drawPolygon(ctx: CanvasRenderingContext2D, particle: PolygonParticle) {
  ctx.beginPath();
  for (let i = 0; i < particle.sides; i += 1) {
    const current = particle.angle + (i * Math.PI * 2) / particle.sides;
    const px = particle.x + Math.cos(current) * particle.radius;
    const py = particle.y + Math.sin(current) * particle.radius;
    if (i === 0) {
      ctx.moveTo(px, py);
    } else {
      ctx.lineTo(px, py);
    }
  }
  ctx.closePath();

  ctx.globalAlpha = particle.alpha;
  ctx.fillStyle = particle.color;
  ctx.fill();
  ctx.globalAlpha = Math.min(0.9, particle.alpha + 0.25);
  ctx.strokeStyle = "#ffffff";
  ctx.lineWidth = 1;
  ctx.stroke();
}

function updateParticle(particle: PolygonParticle) {
  let speedFactor = 1;
  let ax = 0;
  let ay = 0;
  const randomDriftX = rand(-0.005, 0.005);
  const randomDriftY = rand(-0.005, 0.005);

  if (mouse.active) {
    const dx = mouse.x - particle.x;
    const dy = mouse.y - particle.y;
    const distance = Math.hypot(dx, dy);

    if (distance < ATTRACT_RADIUS) {
      const ratio = 1 - distance / ATTRACT_RADIUS;
      speedFactor = 1 + ratio * 0.55;
      if (distance > 0.001) {
        ax = (dx / distance) * ratio * 0.012;
        ay = (dy / distance) * ratio * 0.012;
      }
    }
  }

  particle.vx = particle.vx * 0.992 + ax + randomDriftX;
  particle.vy = particle.vy * 0.992 + ay + randomDriftY;

  const maxSpeed = 1.2;
  particle.vx = Math.max(-maxSpeed, Math.min(maxSpeed, particle.vx));
  particle.vy = Math.max(-maxSpeed, Math.min(maxSpeed, particle.vy));

  particle.x += particle.vx * speedFactor;
  particle.y += particle.vy * speedFactor;
  particle.angle += particle.rotateSpeed * speedFactor;

  if (particle.x < -40) particle.x = width + 40;
  if (particle.x > width + 40) particle.x = -40;
  if (particle.y < -40) particle.y = height + 40;
  if (particle.y > height + 40) particle.y = -40;
}

function animate() {
  const canvas = canvasRef.value;
  if (!canvas) {
    return;
  }

  const ctx = canvas.getContext("2d");
  if (!ctx) {
    return;
  }

  ctx.clearRect(0, 0, width, height);

  for (const particle of particles) {
    updateParticle(particle);
    drawPolygon(ctx, particle);
  }

  frameId = window.requestAnimationFrame(animate);
}

function onMouseMove(event: MouseEvent) {
  const canvas = canvasRef.value;
  if (!canvas) {
    return;
  }
  const rect = canvas.getBoundingClientRect();
  mouse.x = event.clientX - rect.left;
  mouse.y = event.clientY - rect.top;
  mouse.active = true;
}

function onMouseLeave() {
  mouse.active = false;
}

onMounted(() => {
  resizeCanvas();
  animate();
  window.addEventListener("resize", resizeCanvas);
  window.addEventListener("mousemove", onMouseMove);
  window.addEventListener("mouseout", onMouseLeave);
});

onBeforeUnmount(() => {
  window.cancelAnimationFrame(frameId);
  window.removeEventListener("resize", resizeCanvas);
  window.removeEventListener("mousemove", onMouseMove);
  window.removeEventListener("mouseout", onMouseLeave);
});
</script>
