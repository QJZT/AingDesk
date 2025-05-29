<template>
  <div class="data-list-view">
    <!-- 筛选区域 -->
    <div class="filter-section">
      <n-space vertical :size="4">
        <n-card title="数据筛选" size="small">
          <n-grid :cols="24" :x-gap="12">
            <n-grid-item :span="4">
              <n-form-item label="直播链接">
                <n-select
                  v-model:value="selectedLiveUrl"
                  :options="liveUrlOptions"
                  placeholder="请选择直播链接"
                  clearable
                />
              </n-form-item>
            </n-grid-item>
            <n-grid-item :span="4">
              <n-form-item label="时间范围">
                <n-date-picker
                  v-model:value="dateRange"
                  type="daterange"
                  clearable
                  style="width: 100%"
                />
              </n-form-item>
            </n-grid-item>
            <n-grid-item :span="4">
              <n-button type="primary" @click="fetchData" style="margin-top: 24px">
                查询
              </n-button>
            </n-grid-item>
          </n-grid>
        </n-card>
      </n-space>
    </div>

    <!-- 数据概览 -->
    <div class="summary-section" style="margin-top: 16px">
      <n-card title="数据概览" size="small">
        <n-grid :cols="6" :x-gap="12">
          <n-grid-item v-for="(item, index) in summaryData" :key="index">
            <n-statistic :label="item.label" :value="item.value">
              <template #suffix>
                <span style="font-size: 14px; margin-left: 4px">{{ item.suffix }}</span>
              </template>
            </n-statistic>
          </n-grid-item>
        </n-grid>
      </n-card>
    </div>

    <!-- 表格区域 -->
    <div class="table-section" style="margin-top: 16px">
      <n-card title="统计数据" size="small">
        <n-data-table
          :columns="columns"
          :data="tableData"
          :pagination="pagination"
          :loading="loading"
          striped
        />
      </n-card>
    </div>

    <!-- 折线图区域 -->
    <div class="chart-section" style="margin-top: 16px">
      <n-card title="数据趋势" size="small">
        <div ref="chartRef" style="width: 100%; height: 400px"></div>
      </n-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive, watch } from 'vue';
import { NSpace, NCard, NGrid, NGridItem, NSelect, NDatePicker, NButton, NFormItem, NDataTable, NStatistic } from 'naive-ui';
import * as echarts from 'echarts';

// 筛选条件
const selectedLiveUrl = ref(null);
const dateRange = ref(null); // 改回日期范围选择
const liveUrlOptions = ref([]);

// 加载状态
const loading = ref(false);

// 表格配置
const columns = [
  { title: '时间', key: 'live_date' },
  { title: '弹幕数量', key: 'comment_event' },
  { title: '礼物数量', key: 'gift_event' },
  { title: '进入直播间', key: 'join_event' },
  { title: '分享数量', key: 'share_event' },
  { title: '关注数量', key: 'follow_event' },
  { title: '点赞数量', key: 'like_event' },
];

const tableData = ref([]);
const pagination = reactive({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 30, 40],
  onChange: (page: number) => {
    pagination.page = page;
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize;
    pagination.page = 1;
  }
});

// 数据概览
const summaryData = ref([
  { label: '总弹幕数', value: 0, suffix: '条' },
  { label: '总礼物数', value: 0, suffix: '个' },
  { label: '总进入人数', value: 0, suffix: '人' },
  { label: '总分享数', value: 0, suffix: '次' },
  { label: '总关注数', value: 0, suffix: '人' },
  { label: '总点赞数', value: 0, suffix: '次' },
]);

// 图表相关
const chartRef = ref(null);
let chart: echarts.ECharts | null = null;

// 获取直播链接列表
const fetchLiveUrls = async () => {
  try {
    const response = await fetch('http://127.0.0.1:7072/get_live_urls');
    const data = await response.json();
    liveUrlOptions.value = data.map(url => ({
      label: url,
      value: url
    }));
  } catch (error) {
    console.error('获取直播链接失败:', error);
  }
};

// 获取统计数据
const fetchData = async () => {
  loading.value = true;
  try {
    const params = new URLSearchParams();
    if (selectedLiveUrl.value) {
      params.append('live_url', selectedLiveUrl.value);
    }
    if (dateRange.value && dateRange.value.length === 2) {
      const startDate = new Date(dateRange.value[0]);
      const endDate = new Date(dateRange.value[1]);
      params.append('start_date', startDate.toISOString().split('T')[0]);
      params.append('end_date', endDate.toISOString().split('T')[0]);
    }

    const response = await fetch(`http://127.0.0.1:7072/get_stats?${params.toString()}`);
    const data = await response.json();
    
    // 更新表格数据
    tableData.value = data;
    
    // 更新数据概览 - 累加所有数据
    if (data.length > 0) {
      updateSummary(data);
    } else {
      resetSummary();
    }
    
    // 更新图表
    updateChart(data);
  } catch (error) {
    console.error('获取统计数据失败:', error);
  } finally {
    loading.value = false;
  }
};

// 更新数据概览 - 累加所有数据
const updateSummary = (data) => {
  const totals = {
    comment_event: 0,
    gift_event: 0,
    join_event: 0,
    share_event: 0,
    follow_event: 0,
    like_event: 0
  };
  
  data.forEach(item => {
    totals.comment_event += item.comment_event;
    totals.gift_event += item.gift_event;
    totals.join_event += item.join_event;
    totals.share_event += item.share_event;
    totals.follow_event += item.follow_event;
    totals.like_event += item.like_event;
  });
  
  summaryData.value[0].value = totals.comment_event;
  summaryData.value[1].value = totals.gift_event;
  summaryData.value[2].value = totals.join_event;
  summaryData.value[3].value = totals.share_event;
  summaryData.value[4].value = totals.follow_event;
  summaryData.value[5].value = totals.like_event;
};

// 重置数据概览
const resetSummary = () => {
  summaryData.value.forEach(item => {
    item.value = 0;
  });
};

// 初始化图表
const initChart = () => {
  if (chartRef.value) {
    chart = echarts.init(chartRef.value);
    
    const option = {
      tooltip: {
        trigger: 'axis'
      },
      legend: {
        data: ['弹幕', '礼物', '进入', '分享', '关注', '点赞']
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: []
      },
      yAxis: {
        type: 'value'
      },
      series: [
        {
          name: '弹幕',
          type: 'line',
          data: []
        },
        {
          name: '礼物',
          type: 'line',
          data: []
        },
        {
          name: '进入',
          type: 'line',
          data: []
        },
        {
          name: '分享',
          type: 'line',
          data: []
        },
        {
          name: '关注',
          type: 'line',
          data: []
        },
        {
          name: '点赞',
          type: 'line',
          data: []
        }
      ]
    };
    
    chart.setOption(option);
    
    // 响应窗口大小变化
    window.addEventListener('resize', () => {
      chart?.resize();
    });
  }
};

// 更新图表数据
const updateChart = (data) => {
  if (!chart) return;
  
  // 按时间排序
  const sortedData = [...data].sort((a, b) => 
    new Date(a.live_date).getTime() - new Date(b.live_date).getTime()
  );
  
  // 格式化时间，去掉年份，只显示月-日 时:分
  const times = sortedData.map(item => {
    const date = new Date(item.live_date);
    const month = date.getMonth() + 1;
    const day = date.getDate();
    const hours = date.getHours();
    const minutes = date.getMinutes();
    return `${month}-${day} ${hours}:${minutes < 10 ? '0' + minutes : minutes}`;
  });
  const commentData = sortedData.map(item => item.comment_event);
  const giftData = sortedData.map(item => item.gift_event);
  const joinData = sortedData.map(item => item.join_event);
  const shareData = sortedData.map(item => item.share_event);
  const followData = sortedData.map(item => item.follow_event);
  const likeData = sortedData.map(item => item.like_event);
  
  chart.setOption({
    xAxis: {
      data: times
    },
    series: [
      { data: commentData },
      { data: giftData },
      { data: joinData },
      { data: shareData },
      { data: followData },
      { data: likeData }
    ]
  });
};

// 组件挂载时初始化
onMounted(() => {
  fetchLiveUrls();
  initChart();
  fetchData(); // 初始加载数据
});

// 监听窗口大小变化，重绘图表
window.addEventListener('resize', () => {
  chart?.resize();
});
</script>

<style scoped lang="scss">
@use "@/assets/base";

.data-list-view {
  padding: 16px;
  
  .filter-section,
  .summary-section,
  .table-section,
  .chart-section {
    margin-bottom: 16px;
  }
}
</style>