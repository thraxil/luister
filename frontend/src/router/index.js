import { createRouter, createWebHashHistory } from 'vue-router'
import Index from '@/components/Index.vue'
import Tags from '@/components/Tags.vue'
import Tag from '@/components/Tag.vue'

export default createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', name: 'Index', component: Index },
    { path: '/t/', name: 'Tags', component: Tags },
    { path: '/t/:tagname/', name: 'Tag', component: Tag }
  ]
})