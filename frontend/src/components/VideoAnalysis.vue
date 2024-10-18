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
    <!-- <section class="comments">
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
    </section> -->

    <v-card class="comments">
      <v-card-title>Comments</v-card-title>
      <v-card-text>
        <v-row class="sort-controls">
          <v-col cols="12" sm="6">
            <v-select
              v-model="sortBy"
              :items="sortByOptions"
              label="Sort by"
              @update:model-value="sortComments"
              dense
              outlined
            ></v-select>
          </v-col>
          <v-col cols="12" sm="6">
            <v-select
              v-model="sortOrder"
              :items="sortOrderOptions"
              label="Order"
              @update:model-value="sortComments"
              dense
              outlined
            ></v-select>
          </v-col>
        </v-row>
        <v-list v-if="sortedComments.length > 0">
          <v-list-item v-for="comment in sortedComments" :key="comment.id">
            <v-list-item-title>{{ comment.text }}</v-list-item-title>
            <v-list-item-subtitle :class="['sentiment', comment.sentiment]">
              Sentiment: {{ comment.sentiment }}
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
        if (sortBy.value === 'sentiment') {
          const sentimentOrder = { positive: 3, neutral: 2, negative: 1 };
          compareA = sentimentOrder[a.sentiment];
          compareB = sentimentOrder[b.sentiment];
        } else if (sortBy.value === 'length') {
          compareA = a.text.length;
          compareB = b.text.length;
        }

        if (sortOrder.value === 'asc') {
          return compareA - compareB;
        } else {
          return compareB - compareA;
        }
      });
    });

    const fetchVideoData = async () => {
      console.log(props);
      if (!props.videoId) return;

      emit('loading-change', true);
      isLoading.value = true;
      console.log(props.videoId);
      try {
        const response = await axios.get(`/api/video?videoId=${props.videoId}`);
        videoStats.value = response.data.videoStats;
        comments.value = response.data.comments;
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

    watch(
      () => props.videoId,
      (newVideoId) => {
        if (newVideoId) {
          console.log(newVideoId);
          fetchVideoData();
        }
      }
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
    };
  },
  // data() {
  //   return {
  //     videoId: '',
  //     videoStats: null,
  //     comments: [],
  //     isLoading: false,
  //     sortBy: 'sentiment',
  //     sortOrder: 'desc',
  //     sortByOptions: [
  //       { text: 'Sentiment', value: 'sentiment' },
  //       { text: 'Comment Length', value: 'length' },
  //     ],
  //     sortOrderOptions: [
  //       { text: 'Ascending', value: 'asc' },
  //       { text: 'Descending', value: 'desc' },
  //     ],
  //   };
  // },
  // watch: {
  //   isLoading(newValue) {
  //     this.$emit('loading-change', newValue);
  //   },
  // },
  // computed: {
  //   overallSentiment() {
  //     if (this.comments.length === 0) return 'N/A';

  //     const sentimentCounts = this.comments.reduce((acc, comment) => {
  //       acc[comment.sentiment] = (acc[comment.sentiment] || 0) + 1;
  //       return acc;
  //     }, {});

  //     const maxSentiment = Object.entries(sentimentCounts).reduce((a, b) => (a[1] > b[1] ? a : b));
  //     return maxSentiment[0];
  //   },

  //   sortedComments() {
  //     return [...this.comments].sort((a, b) => {
  //       let compareA, compareB;
  //       if (this.sortBy === 'sentiment') {
  //         const sentimentOrder = { positive: 3, neutral: 2, negative: 1 };
  //         compareA = sentimentOrder[a.sentiment];
  //         compareB = sentimentOrder[b.sentiment];
  //       } else if (this.sortBy === 'length') {
  //         compareA = a.text.length;
  //         compareB = b.text.length;
  //       }

  //       if (this.sortOrder === 'asc') {
  //         return compareA - compareB;
  //       } else {
  //         return compareB - compareA;
  //       }
  //     });
  //   },
  // },
  // methods: {
  //   async fetchVideoData() {
  //     this.isLoading = true;
  //     try {
  //       const response = await axios.get(`/api/video?videoId=${this.videoId}`);
  //       this.videoStats = response.data.videoStats;
  //       this.comments = response.data.comments;
  //     } catch (error) {
  //       console.error('Error fetching video data:', error);
  //       // Handle error (e.g., show error message to user)
  //     } finally {
  //       this.isLoading = false;
  //     }
  //   },
  //   formatNumber(num) {
  //     return new Intl.NumberFormat().format(num);
  //   },
  // },
  // mounted() {
  //   // You might want to trigger this based on user input or a prop
  //   this.videoId = 'example-video-id';
  //   this.fetchVideoData();
  // },
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
