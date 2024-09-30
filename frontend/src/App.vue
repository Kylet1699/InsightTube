<template>
  <div id="app">
    <div class="title-container">
      <h1 class="title"><span class="insight">Insight</span><span class="tube">Tube</span></h1>
      <input v-model="videoURL" placeholder="Enter YouTube Video URL" />
      <button @click="analyzeVideo">Analyze Video</button>
    </div>
    <VideoAnalysis ref="videoAnalysis" />
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
</style>
