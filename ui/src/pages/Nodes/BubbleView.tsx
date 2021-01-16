import { Bubble } from '@nivo/circle-packing'

const defaultRoot = {
  name: '节点',
  children: [
    {
      name: '添加节点以观察状态',
      num: 1,
    },
  ],
}

const BubbleView = (props: any) => (
  <Bubble
    {...props}
    colors={{ scheme: 'blues' }}
    root={
      props.root.length ? { ...defaultRoot, children: props.root.map((d: any) => ({ ...d, num: 1 })) } : defaultRoot
    }
    margin={{ top: 30, right: 30, bottom: 30, left: 30 }}
    padding={15}
    borderWidth={3}
    identity="name"
    value="num"
    isZoomable={false}
  />
)

export default BubbleView
