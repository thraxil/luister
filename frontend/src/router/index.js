import Vue from 'vue'
import Router from 'vue-router'
import Index from '@/components/Index'
import Tags from '@/components/Tags'
import Tag from '@/components/Tag'

Vue.use(Router)

export default new Router({
    routes: [
        { path: '/', name: 'Index', component: Index },
        { path: '/t/', name: 'Tags', component: Tags },
        { path: '/t/:tagname/', name: 'Tag', component: Tag }
    ]
})
