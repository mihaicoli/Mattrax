import test from 'ava'
import { mount } from '@vue/test-utils'
import Header from '@/components/Header.vue'

test('is a Vue instance', (t) => {
  const wrapper = mount(Header)
  t.truthy(wrapper.vm)
})
