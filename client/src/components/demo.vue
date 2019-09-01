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
            namee: 'Cart1321',
            show: '',
            status: ''
        }
    },
    mounted() {
        axios.get("http://localhost:8082" + '/uList')
            .then((response) => {
                console.log(response.data);
                this.status = response.status;
                this.show = response.data;
            })
            .catch((error) => {
                console.log(error);
            });
    }
}
</script>
