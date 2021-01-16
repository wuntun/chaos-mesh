export interface ExperimentBasic {
  name: string
  namespace: string
  labels: object | string[]
  annotations: object | string[]
}

export interface ExperimentTargetPod {
  action: 'pod-failure' | 'pod-kill' | 'container-kill'
  container_name?: string
}

export interface ExperimentScope {
  namespace_selectors: string[]
  label_selectors: object | string[]
  annotation_selectors: object | string[]
  phase_selectors: string[]
  mode: string
  value: string
  pods: Record<string, string[]> | string[]
}

export interface ExperimentTargetNetworkLoss {
  loss: string
  correlation: string
}

export interface ExperimentTargetNetworkDelay {
  latency: string
  jitter: string
  correlation: string
}

export interface ExperimentTargetNetworkDuplicate {
  duplicate: string
  correlation: string
}

export interface ExperimentTargetNetworkCorrupt {
  corrupt: string
  correlation: string
}

export interface ExperimentTargetNetworkBandwidth {
  rate: string
  limit: number
  buffer: number
  minburst: number
  peakrate: number
}

export interface ExperimentTargetNetwork {
  action: 'partition' | 'loss' | 'delay' | 'duplicate' | 'corrupt' | 'bandwidth'
  loss: ExperimentTargetNetworkLoss
  delay: ExperimentTargetNetworkDelay
  duplicate: ExperimentTargetNetworkDuplicate
  corrupt: ExperimentTargetNetworkCorrupt
  bandwidth: ExperimentTargetNetworkBandwidth
  direction: 'from' | 'to' | 'both' | ''
  target_scope?: ExperimentScope
}

export interface ExperimentTargetIO {
  action: 'latency' | 'fault' | 'attrOverride'
  delay?: string
  errno?: number
  attr?: object | string[]
  volume_path: string
  path: string
  percent: number
  methods: string[]
}

export interface CallchainFrame {
  funcname: string
  parameters: string
  predicate: string
}

export interface FailKernelReq {
  callchain: CallchainFrame[]
  failtype: number
  headers: string[]
  probability: number
  times: number
}

export interface ExperimentTargetKernel {
  fail_kern_request: FailKernelReq
}

export interface ExperimentTargetTime {
  time_offset: string
  clock_ids: string[]
  container_names: string[]
}

export interface ExperimentTargetStress {
  stressors: {
    cpu?: {
      workers: number
      load: number
      options: string[]
    }
    memory?: {
      workers: number
      options: string[]
    }
  }
  stressng_stressors: string
  container_name: string
}

export const kind = [
  'PodChaos',
  'NetworkChaos',
  'IoChaos',
  'KernelChaos',
  'TimeChaos',
  'StressChaos',
  'DNSChaos',
] as const
export type ExperimentKind = typeof kind[number]

export interface ExperimentTarget {
  kind: ExperimentKind
  pod_chaos: ExperimentTargetPod
  network_chaos: ExperimentTargetNetwork
  io_chaos: ExperimentTargetIO
  kernel_chaos: ExperimentTargetKernel
  time_chaos: ExperimentTargetTime
  stress_chaos: ExperimentTargetStress
}

export interface ExperimentSchedule {
  cron: string
  duration: string
}

export interface Experiment extends ExperimentBasic {
  scope: ExperimentScope
  target: ExperimentTarget
  scheduler: ExperimentSchedule
}

export interface PhysicExperiment {
  name: string
  namespace?: unknown
  labels?: unknown
  annotations?: unknown
  scope?: unknown
  target: ExperimentTarget
  scheduler: ExperimentSchedule
}
