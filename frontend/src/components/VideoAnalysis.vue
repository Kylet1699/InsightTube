<!-- VideoAnalysis.vue -->
<template>
  <div class="video-analysis">
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

    <v-card class="comments">
      <!-- <v-card-title>Comments</v-card-title> -->
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
</style>
