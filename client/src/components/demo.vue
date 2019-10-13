<template>
<div>
    <p>{{namee}}</p>
    <p>请求状态：{{status}}</p>
    <p>请求之后的数据：{{show}}</p>
</div>
</template>

<script>
import axios from 'axios'
export default {
    data() {
        return {
            namee: 'UserList',
            show: '',
            status: ''
        }
    },
    mounted() {
        axios.get("/api/user/getUserList").then((response) => {
            console.log(response.data);
            this.status = response.status;
            this.show = response.data;
        }).catch((error) => {
            console.log(error);
        });
    }
}
</script>
