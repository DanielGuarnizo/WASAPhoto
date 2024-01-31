import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import HomeView from '../views/HomeView.vue'
import StreamView from '../views/StreamView.vue'
import ProfileView from '../views/ProfileView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: HomeView},
		{path: '/session', component: LoginView},
		{path: '/users/:userid/stream', component: StreamView},
		{path: '/users/:id/profile', component: ProfileView}
	]
})

export default router
