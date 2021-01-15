import { Bubble } from '@nivo/circle-packing'

const root = {
  name: 'nodes',
  children: [
    {
      name: 'k8s-cluster-1',
      num: 1,
    },
    {
      name: 'k8s-cluster-2',
      num: 1,
    },
    {
      name: 'physic-node-1',
      num: 1,
    },
    {
      name: 'physic-node-2',
      num: 1,
    },
  ],
}

const BubbleView = (props: any) => (
  <Bubble
    {...props}
    root={root}
    margin={{ top: 30, right: 30, bottom: 30, left: 30 }}
    padding={15}
    identity="name"
    value="num"
    isZoomable={false}
  />
)

export default BubbleView
