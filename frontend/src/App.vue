<template>
  <v-app class="custom-bg">
    <v-main>
      <div class="title-container">
        <h1 class="title"><span class="insight">Insight</span><span class="tube">Tube</span></h1>
        <v-text-field class="video-url-input" label="Video URL or ID" v-model="videoURL" placeholder="Type here..." />
        <v-btn @click="analyzeVideo" elevation="2" :loading="isLoading">Analyze</v-btn>
      </div>
      <VideoAnalysis ref="videoAnalysis" @loading-change="updateIsLoading" :videoId="videoId" />
    </v-main>
  </v-app>
</template>
<script>
import { ref, watch } from 'vue';
import VideoAnalysis from './components/VideoAnalysis.vue';

export default {
  name: 'App',
  components: {
    VideoAnalysis,
  },
  setup() {
    const videoURL = ref('');
    const videoId = ref('');
    const isLoading = ref(false);
    const videoAnalysis = ref(null);

    const analyzeVideo = () => {
      if (videoURL.value.includes('youtube.com/watch?v=')) {
        videoId.value = videoURL.value.split('v=')[1].split('&')[0]; // Extract ID from URL
      } else if (videoURL.value.includes('youtu.be/')) {
        videoId.value = videoURL.value.split('youtu.be/')[1]; // Extract ID from short URL
      } else {
        videoId.value = videoURL.value; // Assume it's already an ID
      }

      if (videoAnalysis.value) {
        videoAnalysis.value.fetchVideoData();
      }
    };

    // Watch for changes in the video data and update the page title
    watch(
      () => videoAnalysis.value?.videoStats,
      (newStats) => {
        if (newStats) {
          document.title = `InsightTube - ${newStats.title}`;
        } else {
          document.title = 'InsightTube';
        }
      },
      { deep: true }
    );

    const updateIsLoading = (loading) => {
      isLoading.value = loading;
    };
    return {
      videoURL,
      videoId,
      isLoading,
      videoAnalysis,
      analyzeVideo,
      updateIsLoading,
    };
  },
};
</script>

<style>
body {
  margin: 0;
  padding: 0;
}

.custom-bg {
  background-color: #2b2a27 !important;
}

.title-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-top: 20px;
}

.title {
  font-size: 2.5em;
  margin-bottom: 20px;
}

.insight {
  color: white;
}

.tube {
  color: red;
}

.video-url-input {
  width: 30rem;
  color: white !important;
}

.video-url-input ::v-deep .v-field__outline {
  color: white !important;
}

.video-url-input ::v-deep .v-label {
  color: white !important;
}

.v-btn {
  margin-top: 20px;
}
</style>
