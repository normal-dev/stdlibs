<script setup>
document.title = 'stdlibs.com - News'
</script>

<template>
  <v-main>
    <v-container>
      <h1 class="text-h4 font-weight-light">
        News
      </h1>
      <h2 class="text-subtitle-1 mt-1 mb-4 text-medium-emphasis">
        New contributions and algorithm updates
      </h2>

      <v-row class="mt-4">
        <v-col
          sm="12"
          md="6">
          <v-card
            flat>
            <v-card-item>
              <v-card-title>
                Go 1.22 released
              </v-card-title>
              <v-card-subtitle>
                Written 02/06/2024
              </v-card-subtitle>
            </v-card-item>

            <v-card-text>
              Go 1.22 fixed the "for" loop semantics, allows to range over
              integers, comes with perfomance improvements and the new
              math/rand/v2 package as well as a net/http server mocker.
              net/http.ServerMux accepts now methods and wildcards.
            </v-card-text>
          </v-card>
        </v-col>
        <v-col
          sm="12"
          md="6">
          <v-card
            flat>
            <v-card-item>
              <v-card-title>
                We added three new contributions
              </v-card-title>
              <v-card-subtitle>
                Written 04/01/2024
              </v-card-subtitle>
            </v-card-item>

            <v-card-text>
              For Golang, we added traefik/traefik and cloudflare/cloudflare with
              net package contributions. moby/moby, canonical/lxd and
              containerd/containerd are using a lot of I/O calls. While
              gopherjs/gopherjs makes a usage of the go/ast package.
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <v-row>
        <v-col
          sm="12"
          md="6">
          <v-card
            flat>
            <v-card-item>
              <v-card-title>
                We added five new contributions
              </v-card-title>
              <v-card-subtitle>
                Written 04/01/2024
              </v-card-subtitle>
            </v-card-item>

            <v-card-text>
              spf13/cobra adds reflect contributions for Golang and
              vercel/next.js with node:http and node:https for Node.js.
              Microsofts microsoft/code imports node:worker_threads, while
              prettier/prettier leverages the file system API. The last addition
              is ueber-go/zap with time, fmt and bytes usage.
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </v-main>
</template>
