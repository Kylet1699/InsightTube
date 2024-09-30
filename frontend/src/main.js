import { createApp } from 'vue';
import App from './App.vue';

import '@fontsource/roboto/300.css';
import '@fontsource/roboto/400.css';
import '@fontsource/roboto/500.css';
import '@fontsource/roboto/700.css';

// Import Vuetify
import { createVuetify } from 'vuetify';
import 'vuetify/styles';
import { aliases, mdi } from 'vuetify/iconsets/mdi';

// Vuetify components and directives
import * as components from 'vuetify/components';
import * as directives from 'vuetify/directives';

// Create the Vuetify instance
const vuetify = createVuetify({
  components,
  directives,
  icons: {
    defaultSet: 'mdi', // Optional icon set
    aliases,
    sets: { mdi },
  },
});

const app = createApp(App);

app.use(vuetify); // Add Vuetify to the app

app.mount('#app');
