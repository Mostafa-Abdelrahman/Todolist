import {createRouter,createWebHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import TodayView from '../views/TodayView.vue'
import UpcomingView from '../views/UpcomingView.vue'
import ListView from '../views/ListView.vue'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {path: '/', component: HomeView},
        {path: '/today', component: TodayView},
        {path: '/upcoming', component: UpcomingView},
        {path: '/list/:id', component: ListView}
    ]
})

export default router