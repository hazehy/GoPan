<template>
  <div class="admin-view-stack">
    <div class="admin-stats-grid">
      <div class="card admin-stat-card" v-for="item in props.stats" :key="item.label">
        <p class="muted">{{ item.label }}</p>
        <p class="admin-stat-value">{{ item.value }}</p>
      </div>
    </div>

    <div class="admin-chart-grid">
      <div class="card admin-chart-card">
        <h3 class="admin-subtitle">用户状态分布</h3>
        <div class="admin-bar-row">
          <span class="admin-bar-label">活跃用户</span>
          <div class="admin-bar-track">
            <div class="admin-bar-fill admin-bar-fill-primary" :style="{ width: `${props.activeUserPercent}%` }"></div>
          </div>
          <span class="admin-bar-value">{{ props.activeUserPercent }}%</span>
        </div>
        <div class="admin-bar-row">
          <span class="admin-bar-label">禁用用户</span>
          <div class="admin-bar-track">
            <div class="admin-bar-fill admin-bar-fill-secondary" :style="{ width: `${props.disabledUserPercent}%` }"></div>
          </div>
          <span class="admin-bar-value">{{ props.disabledUserPercent }}%</span>
        </div>
      </div>

      <div class="card admin-chart-card">
        <h3 class="admin-subtitle">文件类型分布（按扩展名）</h3>
        <div class="admin-pie-layout" v-if="props.extDistribution.length">
          <div class="admin-pie-chart" :style="{ background: props.extPieGradient }"></div>
          <div class="admin-pie-legend">
            <div class="admin-pie-legend-item" v-for="item in props.extDistribution" :key="item.ext">
              <span class="admin-pie-dot" :style="{ background: item.color }"></span>
              <span class="admin-cell-ellipsis admin-tooltip" :data-tooltip="`${item.ext} ${item.percent}%`" :title="`${item.ext} ${item.percent}%`">
                {{ item.ext }} {{ item.percent }}%
              </span>
            </div>
          </div>
        </div>
        <p class="muted" v-if="!props.extDistribution.length">暂无文件类型数据</p>
      </div>
    </div>

    <div class="card admin-section">
      <h3 class="admin-subtitle">数据概览表</h3>
      <table class="table admin-overview-table">
        <thead>
          <tr>
            <th>指标</th>
            <th>数值</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in props.stats" :key="`table-${item.label}`">
            <td>{{ item.label }}</td>
            <td>{{ item.value }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  stats: Array<{ label: string; value: string | number }>;
  activeUserPercent: number;
  disabledUserPercent: number;
  extDistribution: Array<{ ext: string; count: number; percent: number; color: string }>;
  extPieGradient: string;
}

const props = defineProps<Props>();
</script>
