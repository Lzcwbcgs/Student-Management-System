// import { defineComponent, ref } from 'vue'
// import { useRouter } from 'vue-router'
// import { ElMessage, ElMessageBox } from 'element-plus'
// import {
//   HomeFilled,
//   Reading,
//   User,
//   UserFilled,
//   School,
//   List,
//   Expand,
//   Fold
// } from '@element-plus/icons-vue'

// export default defineComponent({
//   components: {
//     HomeFilled,
//     Reading,
//     User,
//     UserFilled,
//     School,
//     List,
//     Expand,
//     Fold
//   },
//   setup() {
//     const isCollapse = ref(false)
//     const router = useRouter()
//     const routes = router.options.routes[1].children

//     const handleCommand = (command) => {
//       if (command === 'logout') {
//         // 添加确认对话框
//         ElMessageBox.confirm(
//           '确定要退出登录吗？',
//           '提示',
//           {
//             confirmButtonText: '确定',
//             cancelButtonText: '取消',
//             type: 'warning',
//           }
//         ).then(() => {
//           // 清除所有相关存储项
//           localStorage.removeItem('token');
//           localStorage.removeItem('userType');
//           localStorage.removeItem('userId');
//           localStorage.removeItem('remembered_username');
//           localStorage.removeItem('remembered_type');
          
//           // 显示退出成功消息
//           ElMessage.success('退出登录成功');
          
//           // 延迟跳转，让用户看到成功消息
//           setTimeout(() => {
//             // 使用router.push而不是window.location.href
//             router.push('/login').then(() => {
//               // 跳转成功后刷新页面确保状态完全重置
//               window.location.reload();
//             });
//           }, 1000);
          
//           // 调试确认
//           console.log('退出登录，已清除的storage:', {
//             token: localStorage.getItem('token'),
//             userType: localStorage.getType('userType')
//           });
//         }).catch(() => {
//           // 用户取消退出
//           console.log('用户取消退出登录');
//         });
//       } else if (command === 'profile') {
//         // 处理个人信息点击
//         ElMessage.info('个人信息功能待开发');
//         // 或者跳转到个人信息页面
//         // router.push('/profile');
//       }
//     };

//     return {
//       isCollapse,
//       routes,
//       handleCommand
//     }
//   }
// })