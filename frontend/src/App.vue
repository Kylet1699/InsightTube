<template>
  <div id="app">
    <div class="title-container">
      <h1 class="title"><span class="insight">Insight</span><span class="tube">Tube</span></h1>
      <v-text-field class="video-url-input" label="Video URL or ID" v-model="videoURL" placeholder="Type here..." />
      <v-btn @click="analyzeVideo" elevation="2" :loading="this.isLoading">Analyze</v-btn>
    </div>
    <VideoAnalysis ref="videoAnalysis" @loading-change="updateIsLoading" />
  </div>
</template>

<script>
import VideoAnalysis from './components/VideoAnalysis.vue';

export default {
  name: 'App',
  components: {
    VideoAnalysis,
  },
  data() {
    return {
      videoURL: '',
      videoId: '',
      isLoading: false,
    };
  },
  methods: {
    analyzeVideo() {
      // ex URL: https://www.youtube.com/watch?v=AjFWyTCDink
      if (this.videoURL.includes('www.youtube.com'))
        this.videoId = this.videoURL.split('v=')[1]; // if user inputs URL, get ID
      else this.videoId = this.videoURL; // if user inputs ID

      this.$refs.videoAnalysis.videoId = this.videoId;
      this.$refs.videoAnalysis.fetchVideoData();
    },
    updateIsLoading(isLoading) {
      console.log(isLoading);
      this.isLoading = isLoading;
    },
  },
};
</script>

<style>
body {
  margin: 0;
  padding: 0;
  background-color: #f0f0f0;
}

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background-color: #2b2a27;
  min-height: 100vh;
}

.title-container {
  display: flex;
  flex-direction: column;
  align-items: center;
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
  color: white;
  border-radius: 3rem;
}
</style>
