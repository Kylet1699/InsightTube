<!-- App.vue -->
<template>
  <div id="app">
    <h1>InsightTube</h1>
    <!-- Input field for YouTube Video ID -->
    <input v-model="videoId" placeholder="Enter YouTube Video ID" />
    <!-- Button to trigger comment analysis -->
    <button @click="fetchComments">Analyze Comments</button>
    <!-- Conditional rendering of comments analysis results -->
    <div v-if="comments.length > 0">
      <h2>Comments Analysis</h2>
      <ul>
        <!-- Loop through comments and display text and sentiment -->
        <li v-for="comment in comments" :key="comment.text">{{ comment.text }} - Sentiment: {{ comment.sentiment }}</li>
      </ul>
    </div>
  </div>
</template>

<script>
// Import axios for making HTTP requests
import axios from 'axios';

export default {
  name: 'App',
  // Define reactive data properties
  data() {
    return {
      videoId: '', // Stores the user-input YouTube Video ID
      comments: [], // Stores the fetched and analyzed comments
    };
  },
  methods: {
    // Method to fetch and analyze comments
    async fetchComments() {
      try {
        // Make GET request to our backend API
        const response = await axios.get(`/api/comments?videoId=${this.videoId}`);
        // Update comments data with the response
        this.comments = response.data;
      } catch (error) {
        // Log any errors that occur during the request
        console.error('Error fetching comments:', error);
        // TODO: Add user-friendly error handling (e.g., display error message to user)
      }
    },
  },
};
</script>

<style>
/* Basic styles for the app */
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
