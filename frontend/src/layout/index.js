import { defineComponent, ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  HomeFilled,
  Reading,
  User,
  UserFilled,
  School,
  List
} from '@element-plus/icons-vue'

export default defineComponent({
  components: {
    HomeFilled,
    Reading,
    User,
    UserFilled,
    School,
    List
  },
  setup() {
    const isCollapse = ref(false)
    const router = useRouter()
    const routes = router.options.routes[1].children

    const handleCommand = (command) => {
      if (command === 'logout') {
        localStorage.removeItem('token')
        router.push('/login')
      }
    }

    return {
      isCollapse,
      routes,
      handleCommand
    }
  }
})
