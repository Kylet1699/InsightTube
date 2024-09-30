<!-- VideoAnalysis.vue -->
<template>
  <div class="video-analysis">
    <!-- Overview Section -->
    <section class="overview">
      <h2>Video Overview</h2>
      <div v-if="videoStats">
        <p><strong>Title:</strong> {{ videoStats.title }}</p>
        <p><strong>View Count:</strong> {{ formatNumber(videoStats.viewCount) }}</p>
        <p><strong>Like Count:</strong> {{ formatNumber(videoStats.likeCount) }}</p>
        <p><strong>Comment Count:</strong> {{ formatNumber(videoStats.commentCount) }}</p>
        <p><strong>Overall Sentiment:</strong> {{ overallSentiment }}</p>
      </div>
      <div v-else>
        <p>No video data available</p>
      </div>
    </section>

    <!-- Comments Section -->
    <section class="comments">
      <h2>Comments</h2>
      <div v-if="comments.length > 0">
        <div v-for="comment in comments" :key="comment.id" class="comment">
          <p>{{ comment.text }}</p>
          <p class="sentiment" :class="comment.sentiment">Sentiment: {{ comment.sentiment }}</p>
        </div>
      </div>
      <div v-else>
        <p>No comments available</p>
      </div>
    </section>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'VideoAnalysis',
  data() {
    return {
      videoId: '',
      videoStats: null,
      comments: [],
      isLoading: false,
    };
  },
  watch: {
    isLoading(newValue) {
      this.$emit('loading-change', newValue);
    },
  },
  computed: {
    overallSentiment() {
      if (this.comments.length === 0) return 'N/A';

      const sentimentCounts = this.comments.reduce((acc, comment) => {
        acc[comment.sentiment] = (acc[comment.sentiment] || 0) + 1;
        return acc;
      }, {});

      const maxSentiment = Object.entries(sentimentCounts).reduce((a, b) => (a[1] > b[1] ? a : b));
      return maxSentiment[0];
    },
  },
  methods: {
    async fetchVideoData() {
      this.isLoading = true;
      try {
        const response = await axios.get(`/api/video?videoId=${this.videoId}`);
        this.videoStats = response.data.videoStats;
        this.comments = response.data.comments;
      } catch (error) {
        console.error('Error fetching video data:', error);
        // Handle error (e.g., show error message to user)
      } finally {
        this.isLoading = false;
      }
    },
    formatNumber(num) {
      return new Intl.NumberFormat().format(num);
    },
  },
  mounted() {
    // You might want to trigger this based on user input or a prop
    this.videoId = 'example-video-id';
    this.fetchVideoData();
  },
};
</script>

<style scoped>
.video-analysis {
  font-family: Arial, sans-serif;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  color: white;
}

.overview,
.comments {
  background-color: #333;
  border-radius: 10px;
  padding: 20px;
  margin-bottom: 20px;
  width: 80%;
  max-width: 800px;
}

.overview {
  border: 2px solid #666;
}

h2 {
  color: #ddd;
  border-bottom: 1px solid #666;
  padding-bottom: 10px;
}

.comment {
  background-color: #393937;
  border-radius: 5px;
  padding: 10px;
  margin-bottom: 10px;
}

.sentiment {
  font-weight: bold;
}

.positive {
  color: #4caf50;
}
.negative {
  color: #f44336;
}
.neutral {
  color: #ffc107;
}
</style>
