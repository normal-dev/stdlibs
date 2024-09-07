<script setup>
document.title = 'stdlibs.com - News'

// YYYY-MM-DD
const news = [
  {
    title: 'We added four new contributions',
    date: '2024-04-15',
    text: `The famous logger prometheus/prometheus has been added to the list of
  Go contributions. Also for Go, rclone/rclone uses the file system API. For
  Node.js we added angular/angular with path, http and fs API calls, as well
  as sveltejs/svelte with assertion usages.`
  },
  {
    title: 'Go 1.22 released',
    date: '2024-02-06',
    text: `Go 1.22 fixed the "for" loop semantics, allows to range over
integers, comes with perfomance improvements and the new
math/rand/v2 package as well as a net/http server mocker.
net/http.ServerMux accepts now methods and wildcards.`
  },
  {
    title: 'We added three new contributions',
    date: '2024-01-02',
    text: `For Golang, we added traefik/traefik and cloudflare/cloudflare
with net package contributions. moby/moby, canonical/lxd and
containerd/containerd are using a lot of I/O calls. While
gopherjs/gopherjs makes a usage of the go/ast package.`
  },
  {
    title: 'We added five new contributions',
    date: '2023-12-13',
    text: `spf13/cobra adds reflect contributions for Golang and
vercel/next.js with node:http and node:https for Node.js.
Microsofts microsoft/code imports node:worker_threads, while
prettier/prettier leverages the file system API. The latest addition
is ueber-go/zap with time, fmt and bytes usages.`
  },
  {
    title: 'Node.js 20.9.0 released',
    date: '2023-10-24',
    text: `Node.js version 20.9.0 marked the switch to Long-Term Support (LTS),
meaning it will receive critical bug fixes for an extended period. This
release also introduced a new crypto.hash() function for faster one-time
hashing of smaller inputs.`
  }
]
</script>

<template>
  <v-main>
    <v-container>
      <v-row>
        <v-col
          v-for="(article, index) in news"
          :key="index"
          cols="12"
          md="4"
          sm="12">
          <v-card flat>
            <v-card-item>
              <v-card-title>{{ article.title }}</v-card-title>
              <v-card-subtitle>Written at {{ article.date }}</v-card-subtitle>
            </v-card-item>

            <v-card-text>
              {{ article.text }}
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </v-main>
</template>
