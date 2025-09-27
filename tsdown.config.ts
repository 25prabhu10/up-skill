import { defineConfig } from 'tsdown'

export default defineConfig([
  {
    dts: false,
    entry: ['./src/**/*.ts'],
    platform: 'neutral',
  },
])
