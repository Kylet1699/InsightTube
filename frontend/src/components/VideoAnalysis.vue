<!-- VideoAnalysis.vue -->
<template>
  <div class="video-analysis">
    <!-- Video overview -->
    <v-card class="overview mb-6">
      <v-card-title>Video Overview</v-card-title>
      <v-card-text v-if="videoStats">
        <v-row>
          <v-col cols="12" md="8">
            <p><strong>Title:</strong> {{ videoStats.title }}</p>
            <p><strong>View Count:</strong> {{ formatNumber(videoStats.viewCount) }}</p>
            <p><strong>Like Count:</strong> {{ formatNumber(videoStats.likeCount) }}</p>
            <p><strong>Comment Count:</strong> {{ formatNumber(videoStats.commentCount) }}</p>
            <p>
              <strong>Overall Sentiment:</strong> <span :class="[overallSentiment]">{{ overallSentiment }} </span>
            </p>
          </v-col>
          <v-col cols="12" md="4" class="d-flex justify-center align-center">
            <v-img :src="videoStats.thumbnailUrl" :aspect-ratio="16 / 9" class="video-thumbnail" cover></v-img>
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-text v-else>
        <p>No video data available</p>
      </v-card-text>
    </v-card>

    <!-- D3 Charts -->
    <v-card class="charts-section mb-6">
      <v-card-title class="charts-title">Data Visualizations</v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12" md="6">
            <v-card class="chart-card">
              <v-card-title class="chart-title">Sentiment Distribution</v-card-title>
              <v-card-text class="chart-container">
                <div id="sentiment-pie-chart"></div>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="12" md="6">
            <v-card class="chart-card">
              <v-card-title class="chart-title">Sentiment Over Time</v-card-title>
              <v-card-text class="chart-container">
                <div id="sentiment-line-chart"></div>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- Comments -->
    <v-card class="comments">
      <v-card-title class="d-flex align-center">
        <span class="mr-4">Comments</span>
        <v-spacer></v-spacer>
        <v-select
          v-model="sortBy"
          :items="sortByOptions"
          label="Sort by"
          @update:model-value="sortComments"
          density="compact"
          hide-details
          class="sorting-select mr-2"
          style="max-width: 190px"
        ></v-select>
        <v-select
          v-model="sortOrder"
          :items="sortOrderOptions"
          label="Order"
          @update:model-value="sortComments"
          density="compact"
          hide-details
          class="sorting-select"
          style="max-width: 150px"
        ></v-select>
      </v-card-title>
      <v-card-text>
        <v-list v-if="sortedComments.length > 0">
          <v-list-item v-for="comment in sortedComments" :key="comment.id">
            <v-list-item-title>{{ comment.text }}</v-list-item-title>
            <v-list-item-subtitle>
              <span :class="['sentiment', comment.sentiment]">Sentiment: {{ comment.sentiment }} </span>

              <span class="updated-at"> &emsp;Updated @ {{ comment.updatedAt }} </span>
            </v-list-item-subtitle>
          </v-list-item>
        </v-list>
        <p v-else>No comments available</p>
      </v-card-text>
    </v-card>
  </div>
</template>

<script>
import { ref, computed, watch } from 'vue';
import axios from 'axios';
import * as d3 from 'd3';

export default {
  name: 'VideoAnalysis',
  emits: ['loading-change'],
  props: {
    videoId: String,
  },
  setup(props, { emit }) {
    const videoStats = ref(null);
    const comments = ref([]);
    const sortBy = ref('sentiment');
    const sortOrder = ref('desc');
    const isLoading = ref(false);

    const sortByOptions = [
      { title: 'Sentiment', value: 'sentiment' },
      { title: 'Comment Length', value: 'length' },
      { title: 'Updated At', value: 'updated at' },
    ];

    const sortOrderOptions = [
      { title: 'Ascending', value: 'asc' },
      { title: 'Descending', value: 'desc' },
    ];

    const overallSentiment = computed(() => {
      if (comments.value.length === 0) return 'N/A';

      const sentimentCounts = comments.value.reduce((acc, comment) => {
        acc[comment.sentiment] = (acc[comment.sentiment] || 0) + 1;
        return acc;
      }, {});

      const maxSentiment = Object.entries(sentimentCounts).reduce((a, b) => (a[1] > b[1] ? a : b));
      return maxSentiment[0];
    });

    const sortedComments = computed(() => {
      return [...comments.value].sort((a, b) => {
        let compareA, compareB;

        switch (sortBy.value) {
          case 'sentiment': {
            const sentimentOrder = { positive: 3, neutral: 2, negative: 1 };
            compareA = sentimentOrder[a.sentiment];
            compareB = sentimentOrder[b.sentiment];
            break;
          }
          case 'length':
            compareA = a.text.length;
            compareB = b.text.length;
            break;
          case 'updated at':
            compareA = new Date(a.updatedAt.replace(/ /, 'T')).getTime();
            compareB = new Date(b.updatedAt.replace(/ /, 'T')).getTime();
            break;
          default:
            console.log('Error sorting');
            return 0;
        }

        return sortOrder.value === 'asc' ? compareA - compareB : compareB - compareA;
      });
    });

    const fetchVideoData = async () => {
      if (!props.videoId) return;

      emit('loading-change', true);
      isLoading.value = true;
      try {
        const response = await axios.get(`/api/video?videoId=${props.videoId}`);
        videoStats.value = response.data.videoStats;
        comments.value = response.data.comments.map((comment) => ({
          ...comment,
          date: new Date(comment.updatedAt),
        }));
      } catch (error) {
        console.error('Error fetching video data:', error);
        // Handle error (e.g., show error message to user)
      } finally {
        isLoading.value = false;
        emit('loading-change', false);
      }
    };

    const formatNumber = (num) => {
      return new Intl.NumberFormat().format(num);
    };

    const sortComments = () => {
      // This method is called when sort options change
      // The actual sorting is done in the sortedComments computed property
    };

    const createSentimentDonutChart = () => {
      const data = [
        { sentiment: 'Positive', value: comments.value.filter((c) => c.sentiment === 'positive').length },
        { sentiment: 'Neutral', value: comments.value.filter((c) => c.sentiment === 'neutral').length },
        { sentiment: 'Negative', value: comments.value.filter((c) => c.sentiment === 'negative').length },
      ];

      const width = 300;
      const height = 300;
      const radius = Math.min(width, height) / 2;

      const color = d3
        .scaleOrdinal()
        .domain(['Positive', 'Neutral', 'Negative'])
        .range(['#4CAF50', '#FFC107', '#F44336']);

      const pie = d3
        .pie()
        .value((d) => d.value)
        .sort(null);

      const arc = d3
        .arc()
        .innerRadius(radius * 0.6) // Inner radius for donut
        .outerRadius(radius);

      // Clear any existing SVG
      d3.select('#sentiment-pie-chart').selectAll('*').remove();

      const svg = d3
        .select('#sentiment-pie-chart')
        .append('svg')
        .attr('width', width)
        .attr('height', height)
        .append('g')
        .attr('transform', `translate(${width / 2},${height / 2})`);

      const arcs = svg.selectAll('arc').data(pie(data)).enter().append('g');

      arcs
        .append('path')
        .attr('d', arc)
        .attr('fill', (d) => color(d.data.sentiment));

      // Add centered legend
      const legend = svg.append('g').attr('text-anchor', 'middle').attr('dominant-baseline', 'middle');

      legend
        .selectAll(null)
        .data(data)
        .enter()
        .append('text')
        .attr('y', (d, i) => i * 20 - 20)
        .attr('fill', (d) => color(d.sentiment))
        .text((d) => `${d.sentiment}: ${d.value}`)
        .style('font-size', '12px')
        .style('font-weight', 'bold')
        .style('text-shadow', '1px 1px 1px rgba(0,0,0,0.5)'); // Add a text shadow for better visibility
    };

    const createSentimentLineChart = () => {
      // Process data
      const commentsByDate = d3.group(comments.value, (d) => d3.timeDay.floor(new Date(d.date)));
      const data = Array.from(commentsByDate, ([date, comments]) => ({
        date,
        positive: comments.filter((c) => c.sentiment === 'positive').length,
        neutral: comments.filter((c) => c.sentiment === 'neutral').length,
        negative: comments.filter((c) => c.sentiment === 'negative').length,
      })).sort((a, b) => a.date - b.date);

      const width = 500;
      const height = 300;
      const margin = { top: 20, right: 100, bottom: 40, left: 50 };

      const x = d3
        .scaleTime()
        .domain(d3.extent(data, (d) => d.date))
        .range([margin.left, width - margin.right]);

      const y = d3
        .scaleLinear()
        .domain([0, d3.max(data, (d) => Math.max(d.positive, d.neutral, d.negative))])
        .nice()
        .range([height - margin.bottom, margin.top]);

      const line = d3
        .line()
        .x((d) => x(d.date))
        .y((d) => y(d.value));

      // Clear any existing SVG
      d3.select('#sentiment-line-chart').selectAll('*').remove();

      const svg = d3.select('#sentiment-line-chart').append('svg').attr('width', width).attr('height', height);

      // Add X axis
      svg
        .append('g')
        .attr('transform', `translate(0,${height - margin.bottom})`)
        .call(
          d3
            .axisBottom(x)
            .ticks(width > 500 ? 10 : 5)
            .tickSizeOuter(0)
        )
        .call((g) => g.select('.domain').attr('stroke', 'white'))
        .call((g) => g.selectAll('.tick line').attr('stroke', 'white'))
        .call((g) => g.selectAll('.tick text').attr('fill', 'white'));

      // Add Y axis
      svg
        .append('g')
        .attr('transform', `translate(${margin.left},0)`)
        .call(d3.axisLeft(y))
        .call((g) => g.select('.domain').remove())
        .call((g) =>
          g
            .selectAll('.tick line')
            .clone()
            .attr('x2', width - margin.left - margin.right)
            .attr('stroke-opacity', 0.1)
            .attr('stroke', 'white')
        )
        .call((g) => g.selectAll('.tick text').attr('fill', 'white'));

      const sentiments = ['positive', 'neutral', 'negative'];
      const colors = ['#4CAF50', '#FFC107', '#F44336'];

      sentiments.forEach((sentiment, i) => {
        svg
          .append('path')
          .datum(data)
          .attr('fill', 'none')
          .attr('stroke', colors[i])
          .attr('stroke-width', 1.5)
          .attr(
            'd',
            line.y((d) => y(d[sentiment]))
          );
      });

      // Add legend to the right side
      const legend = svg
        .append('g')
        .attr('font-family', 'sans-serif')
        .attr('font-size', 10)
        .attr('text-anchor', 'start')
        .selectAll('g')
        .data(sentiments)
        .join('g')
        .attr('transform', (d, i) => `translate(${width - margin.right + 10},${i * 20 + 20})`);

      legend
        .append('circle')
        .attr('cx', 15)
        .attr('cy', 10)
        .attr('r', 6)
        .attr('fill', (d, i) => colors[i]);

      legend
        .append('text')
        .attr('x', 24)
        .attr('y', 9.5)
        .attr('dy', '0.35em')
        .text((d) => d.charAt(0).toUpperCase() + d.slice(1))
        .attr('fill', 'white');
    };

    watch(
      () => props.videoId,
      (newVideoId) => {
        if (newVideoId) {
          console.log(newVideoId);
          fetchVideoData();
        }
      }
    );

    watch(
      () => comments.value,
      () => {
        if (comments.value.length > 0) {
          createSentimentLineChart(), createSentimentDonutChart();
        }
      },
      { deep: true }
    );

    return {
      videoStats,
      sortedComments,
      overallSentiment,
      sortBy,
      sortOrder,
      sortByOptions,
      sortOrderOptions,
      formatNumber,
      sortComments,
      isLoading,
      fetchVideoData,
      comments,
    };
  },
};
</script>

<style scoped>
.video-analysis {
  min-height: 100vh;
  color: white;
  margin: 3rem;
  padding: 2rem;
  border-radius: 1rem;
}

.overview {
  background-color: #282c34 !important;
  border-radius: 0.2rem;
  border: #252525 2px solid;
}

.overview p {
  padding-bottom: 0.8rem;
}

.comments {
  margin-top: 1rem;
  background-color: #333 !important;
}

.v-card-title {
  color: #ddd !important;
  border-bottom: 1px solid #666;
}

.v-card-text {
  margin-top: 1rem;
  color: white !important;
}

.sentiment {
  font-weight: bold;
}

.positive {
  color: #4caf50 !important;
}
.negative {
  color: #f44336 !important;
}
.neutral {
  color: #ffc107 !important;
}

.updated-at {
  font-size: 0.7rem;
}

:deep(.v-list) {
  background-color: transparent !important;
}

:deep(.v-list-item) {
  color: white !important;
}

:deep(.v-field__outline) {
  --v-field-border-opacity: 1 !important;
  color: rgba(255, 255, 255, 0.7) !important;
}

:deep(.v-field__input) {
  color: white !important;
}

:deep(.v-label) {
  color: rgba(255, 255, 255, 0.7) !important;
}

:deep(.v-select__selection) {
  color: white !important;
}

/* d3 charts */
.charts-section {
  background-color: #333 !important;
}

.charts-title {
  color: white !important;
  font-size: 1.5em !important;
  padding: 16px;
  background-color: rgba(0, 0, 0, 0.2);
}

.chart-card {
  background-color: #3f3f3f !important;
  border-radius: 8px;
  overflow: hidden;
  height: 100%; /* Ensure both cards have the same height */
}

.chart-title {
  color: white !important;
  text-align: center;
  padding: 16px;
  font-size: 1.2em;
  background-color: rgba(0, 0, 0, 0.2);
}

#sentiment-pie-chart,
#sentiment-line-chart {
  background-color: #3f3f3f;
  border-radius: 8px;
  padding: 16px;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
}

#sentiment-pie-chart text,
#sentiment-line-chart text {
  fill: white;
}

#sentiment-line-chart path {
  stroke-width: 2;
}

#sentiment-line-chart .axis path,
#sentiment-line-chart .axis line {
  stroke: #888;
}

#sentiment-line-chart .axis text {
  font-size: 10px;
}

/* Ensure responsiveness */
@media (max-width: 960px) {
  .chart-card {
    margin-bottom: 20px;
  }
}
</style>
