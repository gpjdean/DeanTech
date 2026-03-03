<template>
  <div class="deployments-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>工作负载</h2>
        <div class="header-actions">
          <el-select v-model="selectedCluster" placeholder="选择集群" style="width: 200px; margin-right: 10px;">
            <el-option
              v-for="cluster in clusters"
              :key="cluster.id"
              :label="cluster.name"
              :value="cluster.id"
            />
          </el-select>
          <el-select v-model="selectedNamespace" placeholder="选择命名空间" style="width: 200px; margin-right: 10px;" :loading="namespacesLoading" filterable>
          <el-option
            v-for="namespace in namespaces"
            :key="namespace"
            :label="namespace"
            :value="namespace"
          />
        </el-select>
          <el-button type="primary" size="small" @click="refreshCurrentWorkload">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
          <el-button type="primary" size="small" @click="createWorkload">
            <el-icon><Plus /></el-icon>
            创建
          </el-button>
        </div>
      </div>
    </el-card>

    <el-card shadow="hover" class="deployments-card">
      <template #header>
        <div class="card-header">
          <el-tabs v-model="activeWorkloadTab" @tab-change="handleTabChange" type="card">
            <el-tab-pane label="无状态负载 (Deployment)" name="deployment">
              <template #label>
                <div class="tab-label">
                  <el-icon><Grid /></el-icon>
                  <span>无状态负载 (Deployment)</span>
                </div>
              </template>
            </el-tab-pane>
            <el-tab-pane label="有状态负载 (StatefulSet)" name="statefulset">
              <template #label>
                <div class="tab-label">
                  <el-icon><Box /></el-icon>
                  <span>有状态负载 (StatefulSet)</span>
                </div>
              </template>
            </el-tab-pane>
            <el-tab-pane label="守护进程集 (DaemonSet)" name="daemonset">
              <template #label>
                <div class="tab-label">
                  <el-icon><Cloudy /></el-icon>
                  <span>守护进程集 (DaemonSet)</span>
                </div>
              </template>
            </el-tab-pane>
            <el-tab-pane label="普通任务 (Job)" name="job">
              <template #label>
                <div class="tab-label">
                  <el-icon><Timer /></el-icon>
                  <span>普通任务 (Job)</span>
                </div>
              </template>
            </el-tab-pane>
            <el-tab-pane label="定时任务 (CronJob)" name="cronjob">
              <template #label>
                <div class="tab-label">
                  <el-icon><Clock /></el-icon>
                  <span>定时任务 (CronJob)</span>
                </div>
              </template>
            </el-tab-pane>
          </el-tabs>
        </div>
      </template>
      
      <!-- 无状态负载 (Deployment) -->
      <div v-if="activeWorkloadTab === 'deployment'">
        <el-table 
          :data="deployments" 
          style="width: 100%; color: #303133; font-size: 14px; line-height: 1.5;"
          stripe
          v-loading="loading.deployment"
          v-if="selectedCluster"
          :row-style="{ height: 'auto', backgroundColor: '#ffffff' }"
          :cell-style="{ 
            'white-space': 'normal', 
            'word-break': 'break-all', 
            'min-height': '40px', 
            'padding': '12px', 
            'color': '#303133', 
            'font-size': '14px',
            'line-height': '1.5'
          }"
          :header-cell-style="{ 
            'color': '#303133', 
            'font-weight': 'bold', 
            'background-color': '#f5f7fa',
            'font-size': '14px'
          }"
          empty-text="No Data"
          :row-key="(row: any) => `${row.namespace}-${row.name}`"
        >
          <el-table-column prop="name" label="名称" min-width="200" show-overflow-tooltip />
          <el-table-column prop="namespace" label="命名空间" width="150" show-overflow-tooltip />
          <el-table-column prop="status" label="状态" width="120">
            <template #default="scope">
              <el-tag :type="getStatusTagType(scope.row.status)" size="medium">
                {{ scope.row.status }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="replicas" label="副本数" width="120">
            <template #header>
              <div class="column-header">
                <span>副本数</span>
                <el-tooltip content="正常Pod/总Pod" placement="top">
                  <el-icon class="info-icon"><InfoFilled /></el-icon>
                </el-tooltip>
              </div>
            </template>
            <template #default="scope">
              {{ scope.row.readyReplicas }}/{{ scope.row.replicas }}
            </template>
          </el-table-column>
          <el-table-column prop="image" label="镜像" min-width="200" show-overflow-tooltip />
          <el-table-column prop="createdAt" label="创建时间" min-width="180" align="center">
            <template #default="scope">
              <span style="font-weight: 500; color: #606266;">
                {{ formatDateTime(scope.row.createdAt) }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="220" fixed="right">
            <template #default="scope">
              <el-space>
                <el-button type="primary" size="small" @click="viewWorkloadDetails('deployment', scope.row)">
                  <el-icon><View /></el-icon>
                  查看
                </el-button>
                <el-dropdown @command="(command: string) => handleWorkloadCommand(command, 'deployment', scope.row)">
                  <el-button size="small">
                    更多
                    <el-icon class="el-icon--right"><ArrowDown /></el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="scale">
                        <el-icon><ScaleToOriginal /></el-icon>
                        扩缩容
                      </el-dropdown-item>
                      <el-dropdown-item command="restart">
                        <el-icon><Refresh /></el-icon>
                        重启
                      </el-dropdown-item>
                      <el-dropdown-item command="versions">
                        <el-icon><Time /></el-icon>
                        版本记录
                      </el-dropdown-item>
                      <el-dropdown-item command="delete" type="danger">
                        <el-icon><Delete /></el-icon>
                        删除
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </el-space>
            </template>
          </el-table-column>
        </el-table>
      </div>
      
      <!-- 有状态负载 (StatefulSet) -->
      <div v-else-if="activeWorkloadTab === 'statefulset'">
        <el-table 
          :data="statefulsets" 
          style="width: 100%; color: #303133; font-size: 14px; line-height: 1.5;"
          stripe
          v-loading="loading.statefulset"
          v-if="selectedCluster"
          :row-style="{ height: 'auto', backgroundColor: '#ffffff' }"
          :cell-style="{ 
            'white-space': 'normal', 
            'word-break': 'break-all', 
            'min-height': '40px', 
            'padding': '12px', 
            'color': '#303133', 
            'font-size': '14px',
            'line-height': '1.5'
          }"
          :header-cell-style="{ 
            'color': '#303133', 
            'font-weight': 'bold', 
            'background-color': '#f5f7fa',
            'font-size': '14px'
          }"
          empty-text="No Data"
          :row-key="(row: any) => `${row.namespace}-${row.name}`"
        >
          <el-table-column prop="name" label="名称" min-width="200" show-overflow-tooltip />
          <el-table-column prop="namespace" label="命名空间" width="150" show-overflow-tooltip />
          <el-table-column prop="status" label="状态" width="120">
            <template #default="scope">
              <el-tag :type="getStatusTagType(scope.row.status)" size="medium">
                {{ scope.row.status }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="replicas" label="副本数" width="120">
            <template #header>
              <div class="column-header">
                <span>副本数</span>
                <el-tooltip content="正常Pod/总Pod" placement="top">
                  <el-icon class="info-icon"><InfoFilled /></el-icon>
                </el-tooltip>
              </div>
            </template>
            <template #default="scope">
              {{ scope.row.readyReplicas }}/{{ scope.row.replicas }}
            </template>
          </el-table-column>
          <el-table-column prop="image" label="镜像" min-width="200" show-overflow-tooltip />
          <el-table-column prop="createdAt" label="创建时间" min-width="180" align="center">
            <template #default="scope">
              <span style="font-weight: 500; color: #606266;">
                {{ formatDateTime(scope.row.createdAt) }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="220" fixed="right">
            <template #default="scope">
              <el-space>
                <el-button type="primary" size="small" @click="viewWorkloadDetails('statefulset', scope.row)">
                  <el-icon><View /></el-icon>
                  查看
                </el-button>
                <el-dropdown @command="(command: string) => handleWorkloadCommand(command, 'statefulset', scope.row)">
                  <el-button size="small">
                    更多
                    <el-icon class="el-icon--right"><ArrowDown /></el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="scale">
                        <el-icon><ScaleToOriginal /></el-icon>
                        扩缩容
                      </el-dropdown-item>
                      <el-dropdown-item command="delete" type="danger">
                        <el-icon><Delete /></el-icon>
                        删除
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </el-space>
            </template>
          </el-table-column>
        </el-table>
      </div>
      
      <!-- 守护进程集 (DaemonSet) -->
      <div v-else-if="activeWorkloadTab === 'daemonset'">
        <el-table 
          :data="daemonsets" 
          style="width: 100%; color: #303133; font-size: 14px; line-height: 1.5;"
          stripe
          v-loading="loading.daemonset"
          v-if="selectedCluster"
          :row-style="{ height: 'auto', backgroundColor: '#ffffff' }"
          :cell-style="{ 
            'white-space': 'normal', 
            'word-break': 'break-all', 
            'min-height': '40px', 
            'padding': '12px', 
            'color': '#303133', 
            'font-size': '14px',
            'line-height': '1.5'
          }"
          :header-cell-style="{ 
            'color': '#303133', 
            'font-weight': 'bold', 
            'background-color': '#f5f7fa',
            'font-size': '14px'
          }"
          empty-text="No Data"
          :row-key="(row: any) => `${row.namespace}-${row.name}`"
        >
          <el-table-column prop="name" label="名称" min-width="200" show-overflow-tooltip />
          <el-table-column prop="namespace" label="命名空间" width="150" show-overflow-tooltip />
          <el-table-column prop="status" label="状态" width="120">
            <template #default="scope">
              <el-tag :type="getStatusTagType(scope.row.status)" size="medium">
                {{ scope.row.status }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column prop="image" label="镜像" min-width="200" show-overflow-tooltip />
          <el-table-column prop="createdAt" label="创建时间" min-width="180" align="center">
            <template #default="scope">
              <span style="font-weight: 500; color: #606266;">
                {{ formatDateTime(scope.row.createdAt) }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="220" fixed="right">
            <template #default="scope">
              <el-space>
                <el-button type="primary" size="small" @click="viewWorkloadDetails('daemonset', scope.row)">
                  <el-icon><View /></el-icon>
                  查看
                </el-button>
                <el-dropdown @command="(command: string) => handleWorkloadCommand(command, 'daemonset', scope.row)">
                  <el-button size="small">
                    更多
                    <el-icon class="el-icon--right"><ArrowDown /></el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="restart">
                        <el-icon><Refresh /></el-icon>
                        重启
                      </el-dropdown-item>
                      <el-dropdown-item command="delete" type="danger">
                        <el-icon><Delete /></el-icon>
                        删除
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </el-space>
            </template>
          </el-table-column>
        </el-table>
      </div>
      
      <!-- 普通任务 (Job) -->
      <div v-else-if="activeWorkloadTab === 'job'">
        <el-table 
          :data="jobs" 
          style="width: 100%; color: #303133; font-size: 14px; line-height: 1.5;"
          stripe
          v-loading="loading.job"
          v-if="selectedCluster"
          :row-style="{ height: 'auto', backgroundColor: '#ffffff' }"
          :cell-style="{ 
            'white-space': 'normal', 
            'word-break': 'break-all', 
            'min-height': '40px', 
            'padding': '12px', 
            'color': '#303133', 
            'font-size': '14px',
            'line-height': '1.5'
          }"
          :header-cell-style="{ 
            'color': '#303133', 
            'font-weight': 'bold', 
            'background-color': '#f5f7fa',
            'font-size': '14px'
          }"
          empty-text="No Data"
          :row-key="(row: any) => `${row.namespace}-${row.name}`"
        >
          <el-table-column prop="name" label="名称" min-width="200" show-overflow-tooltip />
          <el-table-column prop="namespace" label="命名空间" width="150" show-overflow-tooltip />
          <el-table-column prop="status" label="状态" width="120">
            <template #default="scope">
              <el-tag :type="getStatusTagType(scope.row.status)" size="medium">
                {{ scope.row.status }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="completions" label="完成数" width="120">
            <template #header>
              <div class="column-header">
                <span>完成数</span>
                <el-tooltip content="正常Pod/总Pod" placement="top">
                  <el-icon class="info-icon"><InfoFilled /></el-icon>
                </el-tooltip>
              </div>
            </template>
            <template #default="scope">
              {{ scope.row.succeeded }}/{{ scope.row.completions || 1 }}
            </template>
          </el-table-column>
          <el-table-column prop="image" label="镜像" min-width="200" show-overflow-tooltip />
          <el-table-column prop="createdAt" label="创建时间" min-width="180" align="center">
            <template #default="scope">
              <span style="font-weight: 500; color: #606266;">
                {{ formatDateTime(scope.row.createdAt) }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="220" fixed="right">
            <template #default="scope">
              <el-space>
                <el-button type="primary" size="small" @click="viewWorkloadDetails('job', scope.row)">
                  <el-icon><View /></el-icon>
                  查看
                </el-button>
                <el-dropdown @command="(command: string) => handleWorkloadCommand(command, 'job', scope.row)">
                  <el-button size="small">
                    更多
                    <el-icon class="el-icon--right"><ArrowDown /></el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="restart">
                        <el-icon><Refresh /></el-icon>
                        重启
                      </el-dropdown-item>
                      <el-dropdown-item command="delete" type="danger">
                        <el-icon><Delete /></el-icon>
                        删除
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </el-space>
            </template>
          </el-table-column>
        </el-table>
      </div>
      
      <!-- 定时任务 (CronJob) -->
      <div v-else-if="activeWorkloadTab === 'cronjob'">
        <el-table 
          :data="cronjobs" 
          style="width: 100%; color: #303133; font-size: 14px; line-height: 1.5;"
          stripe
          v-loading="loading.cronjob"
          v-if="selectedCluster"
          :row-style="{ height: 'auto', backgroundColor: '#ffffff' }"
          :cell-style="{ 
            'white-space': 'normal', 
            'word-break': 'break-all', 
            'min-height': '40px', 
            'padding': '12px', 
            'color': '#303133', 
            'font-size': '14px',
            'line-height': '1.5'
          }"
          :header-cell-style="{ 
            'color': '#303133', 
            'font-weight': 'bold', 
            'background-color': '#f5f7fa',
            'font-size': '14px'
          }"
          empty-text="No Data"
          :row-key="(row: any) => `${row.namespace}-${row.name}`"
        >
          <el-table-column prop="name" label="名称" min-width="200" show-overflow-tooltip />
          <el-table-column prop="namespace" label="命名空间" width="150" show-overflow-tooltip />
          <el-table-column prop="status" label="状态" width="120">
            <template #default="scope">
              <el-tag :type="getStatusTagType(scope.row.status)" size="medium">
                {{ scope.row.status }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="schedule" label="调度规则" min-width="150" show-overflow-tooltip />
          <el-table-column prop="lastSchedule" label="上次调度" width="120" align="center">
            <template #default="scope">
              <span style="font-weight: 500; color: #606266;">
                {{ scope.row.lastSchedule ? formatAge(scope.row.lastSchedule) : '-' }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="image" label="镜像" min-width="200" show-overflow-tooltip />
          <el-table-column prop="createdAt" label="创建时间" min-width="180" align="center">
            <template #default="scope">
              <span style="font-weight: 500; color: #606266;">
                {{ formatDateTime(scope.row.createdAt) }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="220" fixed="right">
            <template #default="scope">
              <el-space>
                <el-button type="primary" size="small" @click="viewWorkloadDetails('cronjob', scope.row)">
                  <el-icon><View /></el-icon>
                  查看
                </el-button>
                <el-dropdown @command="(command: string) => handleWorkloadCommand(command, 'cronjob', scope.row)">
                  <el-button size="small">
                    更多
                    <el-icon class="el-icon--right"><ArrowDown /></el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="run">
                        <el-icon><CirclePlus /></el-icon>
                        立即运行
                      </el-dropdown-item>
                      <el-dropdown-item command="delete" type="danger">
                        <el-icon><Delete /></el-icon>
                        删除
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </el-space>
            </template>
          </el-table-column>
        </el-table>
      </div>
      
      <div v-else class="no-cluster-selected">
        <el-empty description="请选择一个集群查看工作负载列表" />
      </div>
    </el-card>
    
    <!-- Deployment详情对话框 -->
    <el-dialog
      v-model="deploymentDetailsVisible"
      title="Deployment详情"
      width="800px"
    >
      <div v-if="selectedDeployment" class="deployment-details">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="Deployment名称">{{ selectedDeployment.name }}</el-descriptions-item>
          <el-descriptions-item label="命名空间">{{ selectedDeployment.namespace }}</el-descriptions-item>
          <el-descriptions-item label="集群">{{ selectedDeployment.cluster }}</el-descriptions-item>
          <el-descriptions-item label="状态">{{ selectedDeployment.status }}</el-descriptions-item>
          <el-descriptions-item label="副本数">{{ selectedDeployment.readyReplicas }}/{{ selectedDeployment.replicas }}</el-descriptions-item>
          <el-descriptions-item label="镜像">{{ selectedDeployment.image }}</el-descriptions-item>
        </el-descriptions>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="deploymentDetailsVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 工作负载扩缩容对话框 -->
    <el-dialog
      v-model="scaleDeploymentVisible"
      title="工作负载扩缩容"
      width="500px"
    >
      <div v-if="selectedDeployment" class="scale-deployment">
        <el-form :model="scaleForm" label-width="120px">
          <el-form-item label="工作负载名称">
            <el-input v-model="selectedDeployment.name" disabled />
          </el-form-item>
          <el-form-item label="当前副本数">
            <el-input v-model="selectedDeployment.replicas" disabled />
          </el-form-item>
          <el-form-item label="目标副本数" prop="replicas">
            <el-input-number v-model="scaleForm.replicas" :min="0" :step="1" style="width: 100%" />
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="scaleDeploymentVisible = false">取消</el-button>
          <el-button type="primary" @click="handleScaleDeployment">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- ReplicaSet版本记录对话框 -->
    <el-dialog
      v-model="replicasetVersionsVisible"
      title="ReplicaSet版本记录"
      width="1000px"
    >
      <div v-if="selectedDeployment" class="replicaset-versions">
        <div class="deployment-info" style="margin-bottom: 20px; padding: 10px; background-color: #f5f7fa; border-radius: 6px;">
          <div style="font-weight: bold; margin-bottom: 8px;">{{ selectedDeployment.name }}</div>
          <div style="color: #606266;">{{ selectedDeployment.namespace }} | {{ selectedDeployment.cluster }}</div>
        </div>
        <el-table 
          :data="replicasets" 
          style="width: 100%"
          stripe
        >
          <el-table-column prop="name" label="名称" min-width="200" show-overflow-tooltip />
          <el-table-column prop="replicas" label="副本数" width="120">
            <template #header>
              <div class="column-header">
                <span>副本数</span>
                <el-tooltip content="正常Pod/总Pod" placement="top">
                  <el-icon class="info-icon"><InfoFilled /></el-icon>
                </el-tooltip>
              </div>
            </template>
            <template #default="scope">
              {{ scope.row.readyReplicas }}/{{ scope.row.replicas }}
            </template>
          </el-table-column>
          <el-table-column prop="image" label="镜像" min-width="250" show-overflow-tooltip />
          <el-table-column prop="createdAt" label="创建时间" min-width="180" align="center">
            <template #default="scope">
              {{ formatDateTime(scope.row.createdAt) }}
            </template>
          </el-table-column>
        </el-table>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="replicasetVersionsVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 创建Deployment对话框 -->
    <el-dialog
      v-model="createDeploymentVisible"
      title="创建无状态负载 (Deployment)"
      width="1200px"
      destroy-on-close
    >
      <!-- 模式切换 -->
      <div style="margin-bottom: 20px;">
        <el-radio-group v-model="createMode" @change="toggleCreateMode">
          <el-radio-button label="form">表单模式</el-radio-button>
          <el-radio-button label="yaml">YAML模式</el-radio-button>
        </el-radio-group>
      </div>

      <!-- 表单模式 -->
      <div v-if="createMode === 'form'" class="create-form">
        <el-form label-width="120px" :model="deploymentForm" size="large">
          <el-tabs v-model="activeFormTab" type="border-card" style="margin-bottom: 20px;">
            <!-- 基本信息 -->
            <el-tab-pane label="基本信息" name="basic">
              <el-form-item label="名称" required>
                <el-input v-model="deploymentForm.name" placeholder="请输入Deployment名称" />
              </el-form-item>

              <el-form-item label="命名空间" required>
                <el-select v-model="deploymentForm.namespace" placeholder="请选择命名空间">
                  <el-option
                    v-for="ns in namespaces"
                    :key="ns"
                    :label="ns"
                    :value="ns"
                  />
                </el-select>
              </el-form-item>

              <el-form-item label="副本数" required>
                <el-input-number v-model="deploymentForm.replicas" :min="1" :step="1" />
              </el-form-item>

              <el-form-item label="镜像" required>
                <el-input v-model="deploymentForm.image" placeholder="请输入镜像地址，例如：nginx:latest" />
              </el-form-item>

              <el-form-item label="容器名称">
                <el-input v-model="deploymentForm.containerName" placeholder="请输入容器名称" />
              </el-form-item>

              <!-- 端口配置 -->
              <el-form-item label="端口配置">
                <div v-for="(port, index) in deploymentForm.ports" :key="index" class="form-group">
                  <div style="display: flex; gap: 10px; align-items: center;">
                    <el-input-number v-model="port.containerPort" :min="1" :max="65535" style="width: 150px;" />
                    <el-select v-model="port.protocol" style="width: 100px;">
                      <el-option label="TCP" value="TCP" />
                      <el-option label="UDP" value="UDP" />
                    </el-select>
                    <el-button type="danger" size="small" @click="deploymentForm.ports.splice(index, 1)">删除</el-button>
                  </div>
                </div>
                <el-button type="primary" size="small" @click="deploymentForm.ports.push({ containerPort: 80, protocol: 'TCP' })">添加端口</el-button>
              </el-form-item>
              
              <!-- 标签配置 -->
              <el-form-item label="标签">
                <div v-for="(label, index) in deploymentForm.labels" :key="index" class="form-group">
                  <div style="display: flex; gap: 10px;">
                    <el-input v-model="label.key" placeholder="键" style="flex: 1;" />
                    <el-input v-model="label.value" placeholder="值" style="flex: 1;" />
                    <el-button type="danger" size="small" @click="deploymentForm.labels.splice(index, 1)">删除</el-button>
                  </div>
                </div>
                <el-button type="primary" size="small" @click="deploymentForm.labels.push({ key: '', value: '' })">添加标签</el-button>
              </el-form-item>

              <!-- 注解配置 -->
              <el-form-item label="注解">
                <div v-for="(annotation, index) in deploymentForm.annotations" :key="index" class="form-group">
                  <div style="display: flex; gap: 10px;">
                    <el-input v-model="annotation.key" placeholder="键" style="flex: 1;" />
                    <el-input v-model="annotation.value" placeholder="值" style="flex: 1;" />
                    <el-button type="danger" size="small" @click="deploymentForm.annotations.splice(index, 1)">删除</el-button>
                  </div>
                </div>
                <el-button type="primary" size="small" @click="deploymentForm.annotations.push({ key: '', value: '' })">添加注解</el-button>
              </el-form-item>
            </el-tab-pane>
            
            <!-- 生命周期 -->
            <el-tab-pane label="生命周期" name="lifecycle">
              <el-card shadow="hover" style="margin-bottom: 20px;">
                <template #header>
                  <div class="card-header">
                    <span>启动探针 (Startup Probe)</span>
                  </div>
                </template>
                <el-form-item label="探针类型">
                  <el-select v-model="deploymentForm.lifecycle.startupProbe.type" placeholder="选择探针类型">
                    <el-option label="HTTP GET" value="httpGet" />
                    <el-option label="TCP Socket" value="tcpSocket" />
                    <el-option label="Exec" value="exec" />
                  </el-select>
                </el-form-item>
                
                <template v-if="deploymentForm.lifecycle.startupProbe.type === 'httpGet'">
                  <el-form-item label="路径">
                    <el-input v-model="deploymentForm.lifecycle.startupProbe.httpGet.path" placeholder="例如：/health" />
                  </el-form-item>
                  <el-form-item label="端口">
                    <el-input-number v-model="deploymentForm.lifecycle.startupProbe.httpGet.port" :min="1" :max="65535" />
                  </el-form-item>
                  <el-form-item label="协议">
                    <el-select v-model="deploymentForm.lifecycle.startupProbe.httpGet.scheme" placeholder="选择协议">
                      <el-option label="HTTP" value="HTTP" />
                      <el-option label="HTTPS" value="HTTPS" />
                    </el-select>
                  </el-form-item>
                </template>
                
                <template v-else-if="deploymentForm.lifecycle.startupProbe.type === 'tcpSocket'">
                  <el-form-item label="端口">
                    <el-input-number v-model="deploymentForm.lifecycle.startupProbe.tcpSocket.port" :min="1" :max="65535" />
                  </el-form-item>
                </template>
                
                <template v-else-if="deploymentForm.lifecycle.startupProbe.type === 'exec'">
                  <el-form-item label="命令">
                    <el-input v-model="deploymentForm.lifecycle.startupProbe.exec.command[0]" placeholder="例如：/bin/sh -c 'echo hello'" />
                  </el-form-item>
                </template>
                
                <el-form-item label="初始延迟时间 (秒)">
                  <el-input-number v-model="deploymentForm.lifecycle.startupProbe.initialDelaySeconds" :min="0" />
                </el-form-item>
                <el-form-item label="检测周期 (秒)">
                  <el-input-number v-model="deploymentForm.lifecycle.startupProbe.periodSeconds" :min="1" />
                </el-form-item>
                <el-form-item label="超时时间 (秒)">
                  <el-input-number v-model="deploymentForm.lifecycle.startupProbe.timeoutSeconds" :min="1" />
                </el-form-item>
                <el-form-item label="成功阈值">
                  <el-input-number v-model="deploymentForm.lifecycle.startupProbe.successThreshold" :min="1" />
                </el-form-item>
                <el-form-item label="失败阈值">
                  <el-input-number v-model="deploymentForm.lifecycle.startupProbe.failureThreshold" :min="1" />
                </el-form-item>
              </el-card>
              
              <el-card shadow="hover" style="margin-bottom: 20px;">
                <template #header>
                  <div class="card-header">
                    <span>存活探针 (Liveness Probe)</span>
                  </div>
                </template>
                <el-form-item label="探针类型">
                  <el-select v-model="deploymentForm.lifecycle.livenessProbe.type" placeholder="选择探针类型">
                    <el-option label="HTTP GET" value="httpGet" />
                    <el-option label="TCP Socket" value="tcpSocket" />
                    <el-option label="Exec" value="exec" />
                  </el-select>
                </el-form-item>
                
                <template v-if="deploymentForm.lifecycle.livenessProbe.type === 'httpGet'">
                  <el-form-item label="路径">
                    <el-input v-model="deploymentForm.lifecycle.livenessProbe.httpGet.path" placeholder="例如：/health" />
                  </el-form-item>
                  <el-form-item label="端口">
                    <el-input-number v-model="deploymentForm.lifecycle.livenessProbe.httpGet.port" :min="1" :max="65535" />
                  </el-form-item>
                  <el-form-item label="协议">
                    <el-select v-model="deploymentForm.lifecycle.livenessProbe.httpGet.scheme" placeholder="选择协议">
                      <el-option label="HTTP" value="HTTP" />
                      <el-option label="HTTPS" value="HTTPS" />
                    </el-select>
                  </el-form-item>
                </template>
                
                <template v-else-if="deploymentForm.lifecycle.livenessProbe.type === 'tcpSocket'">
                  <el-form-item label="端口">
                    <el-input-number v-model="deploymentForm.lifecycle.livenessProbe.tcpSocket.port" :min="1" :max="65535" />
                  </el-form-item>
                </template>
                
                <template v-else-if="deploymentForm.lifecycle.livenessProbe.type === 'exec'">
                  <el-form-item label="命令">
                    <el-input v-model="deploymentForm.lifecycle.livenessProbe.exec.command[0]" placeholder="例如：/bin/sh -c 'echo hello'" />
                  </el-form-item>
                </template>
                
                <el-form-item label="初始延迟时间 (秒)">
                  <el-input-number v-model="deploymentForm.lifecycle.livenessProbe.initialDelaySeconds" :min="0" />
                </el-form-item>
                <el-form-item label="检测周期 (秒)">
                  <el-input-number v-model="deploymentForm.lifecycle.livenessProbe.periodSeconds" :min="1" />
                </el-form-item>
                <el-form-item label="超时时间 (秒)">
                  <el-input-number v-model="deploymentForm.lifecycle.livenessProbe.timeoutSeconds" :min="1" />
                </el-form-item>
                <el-form-item label="成功阈值">
                  <el-input-number v-model="deploymentForm.lifecycle.livenessProbe.successThreshold" :min="1" />
                </el-form-item>
                <el-form-item label="失败阈值">
                  <el-input-number v-model="deploymentForm.lifecycle.livenessProbe.failureThreshold" :min="1" />
                </el-form-item>
              </el-card>
              
              <el-card shadow="hover">
                <template #header>
                  <div class="card-header">
                    <span>就绪探针 (Readiness Probe)</span>
                  </div>
                </template>
                <el-form-item label="探针类型">
                  <el-select v-model="deploymentForm.lifecycle.readinessProbe.type" placeholder="选择探针类型">
                    <el-option label="HTTP GET" value="httpGet" />
                    <el-option label="TCP Socket" value="tcpSocket" />
                    <el-option label="Exec" value="exec" />
                  </el-select>
                </el-form-item>
                
                <template v-if="deploymentForm.lifecycle.readinessProbe.type === 'httpGet'">
                  <el-form-item label="路径">
                    <el-input v-model="deploymentForm.lifecycle.readinessProbe.httpGet.path" placeholder="例如：/health" />
                  </el-form-item>
                  <el-form-item label="端口">
                    <el-input-number v-model="deploymentForm.lifecycle.readinessProbe.httpGet.port" :min="1" :max="65535" />
                  </el-form-item>
                  <el-form-item label="协议">
                    <el-select v-model="deploymentForm.lifecycle.readinessProbe.httpGet.scheme" placeholder="选择协议">
                      <el-option label="HTTP" value="HTTP" />
                      <el-option label="HTTPS" value="HTTPS" />
                    </el-select>
                  </el-form-item>
                </template>
                
                <template v-else-if="deploymentForm.lifecycle.readinessProbe.type === 'tcpSocket'">
                  <el-form-item label="端口">
                    <el-input-number v-model="deploymentForm.lifecycle.readinessProbe.tcpSocket.port" :min="1" :max="65535" />
                  </el-form-item>
                </template>
                
                <template v-else-if="deploymentForm.lifecycle.readinessProbe.type === 'exec'">
                  <el-form-item label="命令">
                    <el-input v-model="deploymentForm.lifecycle.readinessProbe.exec.command[0]" placeholder="例如：/bin/sh -c 'echo hello'" />
                  </el-form-item>
                </template>
                
                <el-form-item label="初始延迟时间 (秒)">
                  <el-input-number v-model="deploymentForm.lifecycle.readinessProbe.initialDelaySeconds" :min="0" />
                </el-form-item>
                <el-form-item label="检测周期 (秒)">
                  <el-input-number v-model="deploymentForm.lifecycle.readinessProbe.periodSeconds" :min="1" />
                </el-form-item>
                <el-form-item label="超时时间 (秒)">
                  <el-input-number v-model="deploymentForm.lifecycle.readinessProbe.timeoutSeconds" :min="1" />
                </el-form-item>
                <el-form-item label="成功阈值">
                  <el-input-number v-model="deploymentForm.lifecycle.readinessProbe.successThreshold" :min="1" />
                </el-form-item>
                <el-form-item label="失败阈值">
                  <el-input-number v-model="deploymentForm.lifecycle.readinessProbe.failureThreshold" :min="1" />
                </el-form-item>
              </el-card>
            </el-tab-pane>
            
            <!-- 环境变量 -->
            <el-tab-pane label="环境变量" name="env">
              <div v-for="(env, index) in deploymentForm.env" :key="index" class="form-group">
                <el-select v-model="env.type" placeholder="选择环境变量类型" style="margin-bottom: 10px;">
                  <el-option label="直接值" value="value" />
                  <el-option label="ConfigMap" value="configMapKeyRef" />
                  <el-option label="Secret" value="secretKeyRef" />
                </el-select>
                <div style="display: flex; gap: 10px; flex-wrap: wrap;">
                  <el-input v-model="env.name" placeholder="变量名" style="flex: 1; min-width: 150px;" />
                  <template v-if="env.type === 'value'">
                    <el-input v-model="env.value" placeholder="变量值" style="flex: 1; min-width: 150px;" />
                  </template>
                  <template v-else-if="env.type === 'configMapKeyRef'">
                    <el-input v-model="env.valueFrom.configMapKeyRef.name" placeholder="ConfigMap名称" style="flex: 1; min-width: 150px;" />
                    <el-input v-model="env.valueFrom.configMapKeyRef.key" placeholder="键" style="flex: 1; min-width: 150px;" />
                  </template>
                  <template v-else-if="env.type === 'secretKeyRef'">
                    <el-input v-model="env.valueFrom.secretKeyRef.name" placeholder="Secret名称" style="flex: 1; min-width: 150px;" />
                    <el-input v-model="env.valueFrom.secretKeyRef.key" placeholder="键" style="flex: 1; min-width: 150px;" />
                  </template>
                  <el-button type="danger" size="small" @click="deploymentForm.env.splice(index, 1)">删除</el-button>
                </div>
              </div>
              <el-button type="primary" size="small" @click="deploymentForm.env.push({ name: '', value: '', valueFrom: { configMapKeyRef: { name: '', key: '' }, secretKeyRef: { name: '', key: '' } }, type: 'value' })">添加环境变量</el-button>
            </el-tab-pane>
            
            <!-- 数据存储 -->
            <el-tab-pane label="数据存储" name="storage">
              <h4 style="margin-top: 0; margin-bottom: 20px;">卷配置</h4>
              <div v-for="(volume, index) in deploymentForm.volumes" :key="index" class="form-group">
                <el-select v-model="volume.type" placeholder="选择卷类型" style="margin-bottom: 10px;">
                  <el-option label="EmptyDir" value="emptyDir" />
                  <el-option label="HostPath" value="hostPath" />
                  <el-option label="ConfigMap" value="configMap" />
                  <el-option label="Secret" value="secret" />
                  <el-option label="PersistentVolumeClaim" value="persistentVolumeClaim" />
                </el-select>
                <div style="display: flex; gap: 10px; flex-wrap: wrap;">
                  <el-input v-model="volume.name" placeholder="卷名称" style="flex: 1; min-width: 150px;" />
                  <template v-if="volume.type === 'emptyDir'">
                    <el-select v-model="volume.emptyDir.medium" placeholder="存储介质">
                      <el-option label="默认" value="" />
                      <el-option label="内存" value="Memory" />
                    </el-select>
                  </template>
                  <template v-else-if="volume.type === 'hostPath'">
                    <el-input v-model="volume.hostPath.path" placeholder="主机路径" style="flex: 1; min-width: 150px;" />
                    <el-select v-model="volume.hostPath.type" placeholder="路径类型">
                      <el-option label="目录" value="Directory" />
                      <el-option label="目录或创建" value="DirectoryOrCreate" />
                      <el-option label="文件" value="File" />
                      <el-option label="文件或创建" value="FileOrCreate" />
                    </el-select>
                  </template>
                  <template v-else-if="volume.type === 'configMap'">
                    <el-input v-model="volume.configMap.name" placeholder="ConfigMap名称" style="flex: 1; min-width: 150px;" />
                  </template>
                  <template v-else-if="volume.type === 'secret'">
                    <el-input v-model="volume.secret.secretName" placeholder="Secret名称" style="flex: 1; min-width: 150px;" />
                  </template>
                  <template v-else-if="volume.type === 'persistentVolumeClaim'">
                    <el-input v-model="volume.persistentVolumeClaim.claimName" placeholder="PVC名称" style="flex: 1; min-width: 150px;" />
                  </template>
                  <el-button type="danger" size="small" @click="deploymentForm.volumes.splice(index, 1)">删除</el-button>
                </div>
              </div>
              <el-button type="primary" size="small" @click="deploymentForm.volumes.push({ name: '', type: 'emptyDir', emptyDir: { medium: '' }, hostPath: { path: '', type: '' }, configMap: { name: '', items: [{ key: '', path: '' }] }, secret: { secretName: '', items: [{ key: '', path: '' }] }, persistentVolumeClaim: { claimName: '' } })">添加卷</el-button>
              
              <h4 style="margin-top: 30px; margin-bottom: 20px;">卷挂载</h4>
              <div v-for="(mount, index) in deploymentForm.volumeMounts" :key="index" class="form-group">
                <div style="display: flex; gap: 10px; flex-wrap: wrap;">
                  <el-select v-model="mount.name" placeholder="选择卷" style="flex: 1; min-width: 150px;">
                    <el-option v-for="volume in deploymentForm.volumes" :key="volume.name" :label="volume.name" :value="volume.name" />
                  </el-select>
                  <el-input v-model="mount.mountPath" placeholder="挂载路径" style="flex: 1; min-width: 150px;" />
                  <el-input v-model="mount.subPath" placeholder="子路径" style="flex: 1; min-width: 150px;" />
                  <el-switch v-model="mount.readOnly" active-text="只读" inactive-text="读写" />
                  <el-button type="danger" size="small" @click="deploymentForm.volumeMounts.splice(index, 1)">删除</el-button>
                </div>
              </div>
              <el-button type="primary" size="small" @click="deploymentForm.volumeMounts.push({ name: '', mountPath: '', subPath: '', readOnly: false })">添加卷挂载</el-button>
            </el-tab-pane>
            
            <!-- Init容器 -->
            <el-tab-pane label="Init容器" name="initContainers">
              <div v-for="(container, index) in deploymentForm.initContainers" :key="index" class="form-group">
                <h5 style="margin-top: 0; margin-bottom: 15px; font-weight: bold;">Init容器 {{ index + 1 }}</h5>
                <el-input v-model="container.name" placeholder="容器名称" style="margin-bottom: 10px;" />
                <el-input v-model="container.image" placeholder="镜像地址" style="margin-bottom: 10px;" />
                <el-input v-model="container.command[0]" placeholder="命令" style="margin-bottom: 10px;" />
                <el-input v-model="container.args[0]" placeholder="参数" style="margin-bottom: 10px;" />
                
                <!-- Init容器资源配置 -->
                <div style="margin-bottom: 15px;">
                  <div style="font-weight: bold; margin-bottom: 5px;">资源配置</div>
                  <div style="display: flex; gap: 10px; margin-bottom: 10px;">
                    <div style="flex: 1;">
                      <div style="font-size: 12px; color: #606266; margin-bottom: 5px;">请求CPU</div>
                      <el-input v-model="container.resources.requests.cpu" placeholder="例如：100m" />
                    </div>
                    <div style="flex: 1;">
                      <div style="font-size: 12px; color: #606266; margin-bottom: 5px;">请求内存</div>
                      <el-input v-model="container.resources.requests.memory" placeholder="例如：128Mi" />
                    </div>
                  </div>
                  <div style="display: flex; gap: 10px;">
                    <div style="flex: 1;">
                      <div style="font-size: 12px; color: #606266; margin-bottom: 5px;">限制CPU</div>
                      <el-input v-model="container.resources.limits.cpu" placeholder="例如：500m" />
                    </div>
                    <div style="flex: 1;">
                      <div style="font-size: 12px; color: #606266; margin-bottom: 5px;">限制内存</div>
                      <el-input v-model="container.resources.limits.memory" placeholder="例如：256Mi" />
                    </div>
                  </div>
                </div>
                
                <el-button type="danger" size="small" @click="deploymentForm.initContainers.splice(index, 1)">删除Init容器</el-button>
              </div>
              <el-button type="primary" size="small" @click="deploymentForm.initContainers.push({ name: '', image: '', command: [''], args: [''], env: [{ name: '', value: '', valueFrom: { configMapKeyRef: { name: '', key: '' }, secretKeyRef: { name: '', key: '' } }, type: 'value' }], resources: { requests: { cpu: '100m', memory: '128Mi' }, limits: { cpu: '500m', memory: '256Mi' } }, volumeMounts: [{ name: '', mountPath: '', subPath: '', readOnly: false }] })">添加Init容器</el-button>
            </el-tab-pane>
            
            <!-- 镜像拉取凭证 -->
            <el-tab-pane label="镜像拉取凭证" name="imagePullSecrets">
              <div v-for="(secret, index) in deploymentForm.imagePullSecrets" :key="index" class="form-group">
                <div style="display: flex; gap: 10px;">
                  <el-input v-model="secret.name" placeholder="Secret名称" style="flex: 1;" />
                  <el-button type="danger" size="small" @click="deploymentForm.imagePullSecrets.splice(index, 1)">删除</el-button>
                </div>
              </div>
              <el-button type="primary" size="small" @click="deploymentForm.imagePullSecrets.push({ name: '' })">添加镜像拉取凭证</el-button>
            </el-tab-pane>
            
            <!-- 升级策略 -->
            <el-tab-pane label="升级策略" name="strategy">
              <el-select v-model="deploymentForm.strategy.type" placeholder="选择升级策略">
                <el-option label="滚动更新" value="RollingUpdate" />
                <el-option label="重建" value="Recreate" />
              </el-select>
              
              <template v-if="deploymentForm.strategy.type === 'RollingUpdate'">
                <el-form-item label="最大可用不可用" style="margin-top: 15px;">
                  <el-input v-model="deploymentForm.strategy.rollingUpdate.maxUnavailable" placeholder="例如：25% 或 1" />
                </el-form-item>
                <el-form-item label="最大额外创建">
                  <el-input v-model="deploymentForm.strategy.rollingUpdate.maxSurge" placeholder="例如：25% 或 1" />
                </el-form-item>
              </template>
            </el-tab-pane>
            
            <!-- 调度策略 -->
            <el-tab-pane label="调度策略" name="scheduling">
              <!-- Node Selector -->
              <h4 style="margin-top: 0; margin-bottom: 20px;">Node Selector</h4>
              <div v-for="(selector, index) in deploymentForm.nodeSelector" :key="index" class="form-group">
                <div style="display: flex; gap: 10px;">
                  <el-input v-model="selector.key" placeholder="键" style="flex: 1;" />
                  <el-input v-model="selector.value" placeholder="值" style="flex: 1;" />
                  <el-button type="danger" size="small" @click="deploymentForm.nodeSelector.splice(index, 1)">删除</el-button>
                </div>
              </div>
              <el-button type="primary" size="small" @click="deploymentForm.nodeSelector.push({ key: '', value: '' })">添加Node Selector</el-button>
              
              <!-- 容忍策略 -->
              <h4 style="margin-top: 30px; margin-bottom: 20px;">容忍策略 (Tolerations)</h4>
              <div v-for="(toleration, index) in deploymentForm.tolerations" :key="index" class="form-group">
                <div style="display: flex; gap: 10px; flex-wrap: wrap;">
                  <el-input v-model="toleration.key" placeholder="键" style="flex: 1; min-width: 150px;" />
                  <el-select v-model="toleration.operator" placeholder="操作符" style="flex: 1; min-width: 150px;">
                    <el-option label="Equal" value="Equal" />
                    <el-option label="Exists" value="Exists" />
                  </el-select>
                  <el-input v-model="toleration.value" placeholder="值" style="flex: 1; min-width: 150px;" />
                  <el-select v-model="toleration.effect" placeholder="影响" style="flex: 1; min-width: 150px;">
                    <el-option label="NoSchedule" value="NoSchedule" />
                    <el-option label="PreferNoSchedule" value="PreferNoSchedule" />
                    <el-option label="NoExecute" value="NoExecute" />
                  </el-select>
                  <el-button type="danger" size="small" @click="deploymentForm.tolerations.splice(index, 1)">删除</el-button>
                </div>
              </div>
              <el-button type="primary" size="small" @click="deploymentForm.tolerations.push({ key: '', operator: 'Equal', value: '', effect: 'NoSchedule' })">添加容忍策略</el-button>
            </el-tab-pane>
            
            <!-- DNS配置 -->
            <el-tab-pane label="DNS配置" name="dns">
              <el-select v-model="deploymentForm.dnsPolicy" placeholder="选择DNS策略">
                <el-option label="ClusterFirst" value="ClusterFirst" />
                <el-option label="Default" value="Default" />
                <el-option label="ClusterFirstWithHostNet" value="ClusterFirstWithHostNet" />
              </el-select>
              
              <h4 style="margin-top: 30px; margin-bottom: 20px;">DNS配置</h4>
              <el-input v-model="deploymentForm.dnsConfig.nameservers[0]" placeholder="DNS服务器，多个用逗号分隔" style="margin-bottom: 10px;" />
              <el-input v-model="deploymentForm.dnsConfig.searches[0]" placeholder="搜索域，多个用逗号分隔" style="margin-bottom: 10px;" />
              <div v-for="(option, index) in deploymentForm.dnsConfig.options" :key="index" class="form-group">
                <div style="display: flex; gap: 10px;">
                  <el-input v-model="option.name" placeholder="选项名称" style="flex: 1;" />
                  <el-input v-model="option.value" placeholder="选项值" style="flex: 1;" />
                  <el-button type="danger" size="small" @click="deploymentForm.dnsConfig.options.splice(index, 1)">删除</el-button>
                </div>
              </div>
              <el-button type="primary" size="small" @click="deploymentForm.dnsConfig.options.push({ name: '', value: '' })">添加DNS选项</el-button>
            </el-tab-pane>
            
            <!-- 其他配置 -->
            <el-tab-pane label="其他配置" name="other">
              <el-form-item label="优雅终止时间 (秒)">
                <el-input-number v-model="deploymentForm.terminationGracePeriodSeconds" :min="0" />
              </el-form-item>
              <el-form-item label="使用主机网络">
                <el-switch v-model="deploymentForm.hostNetwork" active-text="启用" inactive-text="禁用" />
              </el-form-item>
              
              <!-- 资源配置 -->
              <h4 style="margin-top: 30px; margin-bottom: 20px;">资源配置</h4>
              <div style="display: flex; gap: 20px; margin-bottom: 20px;">
                <div style="flex: 1;">
                  <div style="font-weight: bold; margin-bottom: 5px;">请求CPU</div>
                  <el-input v-model="deploymentForm.resources.requests.cpu" placeholder="例如：100m" />
                </div>
                <div style="flex: 1;">
                  <div style="font-weight: bold; margin-bottom: 5px;">请求内存</div>
                  <el-input v-model="deploymentForm.resources.requests.memory" placeholder="例如：128Mi" />
                </div>
              </div>
              <div style="display: flex; gap: 20px;">
                <div style="flex: 1;">
                  <div style="font-weight: bold; margin-bottom: 5px;">限制CPU</div>
                  <el-input v-model="deploymentForm.resources.limits.cpu" placeholder="例如：500m" />
                </div>
                <div style="flex: 1;">
                  <div style="font-weight: bold; margin-bottom: 5px;">限制内存</div>
                  <el-input v-model="deploymentForm.resources.limits.memory" placeholder="例如：256Mi" />
                </div>
              </div>
            </el-tab-pane>
          </el-tabs>
        </el-form>
      </div>

      <!-- YAML模式 -->
      <div v-else class="yaml-mode">
        <!-- 验证结果 -->
        <div v-if="validationResult" class="validation-result" :class="validationResult.includes('✅') ? 'success' : 'error'">
          {{ validationResult }}
        </div>
        <div style="margin-bottom: 10px; display: flex; justify-content: flex-end; gap: 10px;">
          <el-tooltip content="预检会通过 dry-run 验证 YAML 是否符合 Kubernetes 规范" placement="bottom">
            <el-button type="primary" size="small" @click="precheckYaml" :loading="validationLoading">
              <el-icon><Check /></el-icon>
              预检
            </el-button>
          </el-tooltip>
          <el-button type="primary" size="small" @click="createDeployment" :loading="createLoading">
            <el-icon><Plus /></el-icon>
            创建
          </el-button>
        </div>
        <el-input
          v-model="yamlContent"
          type="textarea"
          :rows="25"
          placeholder="请输入YAML内容"
          style="font-family: 'Courier New', Courier, monospace; font-size: 14px;"
        />
      </div>

      <template #footer>
        <span class="dialog-footer" v-if="createMode === 'form'">
          <el-button @click="createDeploymentVisible = false">取消</el-button>
          <el-button type="primary" @click="createDeployment" :loading="createLoading">
            <el-icon><Plus /></el-icon>
            创建
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { Refresh, View, ScaleToOriginal, Delete, ArrowDown, Plus, Timer, Clock, CirclePlus, Cloudy, Check, InfoFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'

// 当前选中的工作负载标签页
const activeWorkloadTab = ref('deployment')

// 加载状态，分别跟踪不同工作负载类型的加载状态
const loading = ref({
  deployment: false,
  statefulset: false,
  daemonset: false,
  job: false,
  cronjob: false
})

// 集群数据类型定义
interface Cluster {
  id: number
  name: string
  apiURL: string
  token: string
  kubeconfig: string
  description: string
  status: string
  version: string
  isActive: boolean
}

// 集群列表
const clusters = ref<Cluster[]>([])

// 命名空间列表
const namespaces = ref<string[]>([])
const namespacesLoading = ref(false)

// 选择的集群和命名空间
const selectedCluster = ref<number | undefined>(undefined)
const selectedNamespace = ref<string>('')

// Deployment数据类型定义
interface Deployment {
  clusterID: number
  name: string
  namespace: string
  cluster: string
  status: string
  replicas: number
  readyReplicas: number
  image: string
  createdAt: string
}

// StatefulSet数据类型定义
interface StatefulSet {
  clusterID: number
  name: string
  namespace: string
  cluster: string
  status: string
  replicas: number
  readyReplicas: number
  image: string
  createdAt: string
}

// DaemonSet数据类型定义
interface DaemonSet {
  clusterID: number
  name: string
  namespace: string
  cluster: string
  status: string
  desiredNumberScheduled: number
  currentNumberScheduled: number
  image: string
  createdAt: string
}

// Job数据类型定义
interface Job {
  clusterID: number
  name: string
  namespace: string
  cluster: string
  status: string
  completions: number
  succeeded: number
  image: string
  createdAt: string
}

// CronJob数据类型定义
interface CronJob {
  clusterID: number
  name: string
  namespace: string
  cluster: string
  status: string
  schedule: string
  lastSchedule: string
  image: string
  createdAt: string
}

// ReplicaSet数据类型定义
interface ReplicaSet {
  name: string
  namespace: string
  replicas: number
  readyReplicas: number
  image: string
  createdAt: string


  deployment: string

}

// 工作负载数据
const deployments = ref<Deployment[]>([])
const statefulsets = ref<StatefulSet[]>([])
const daemonsets = ref<DaemonSet[]>([])
const jobs = ref<Job[]>([])
const cronjobs = ref<CronJob[]>([])

// 对话框相关
const deploymentDetailsVisible = ref(false)
const scaleDeploymentVisible = ref(false)
const replicasetVersionsVisible = ref(false)
const createDeploymentVisible = ref(false)
const selectedDeployment = ref<Deployment | null>(null)
const replicasets = ref<ReplicaSet[]>([])
const scaleForm = ref({
  replicas: 0
})

// 创建Deployment相关
const createMode = ref<'form' | 'yaml'>('form')
const activeFormTab = ref('basic')
const yamlContent = ref('')
const createLoading = ref(false)
const validationResult = ref<string | null>(null)
const validationLoading = ref(false)

// Deployment表单数据
const deploymentForm = ref({
  name: '',
  namespace: '',
  replicas: 1,
  image: '',
  containerName: 'app',
  ports: [{
    containerPort: 80,
    protocol: 'TCP'
  }],
  // 生命周期配置
  lifecycle: {
    startupProbe: {
      type: 'httpGet',
      httpGet: {
        path: '',
        port: 80,
        scheme: 'HTTP'
      },
      tcpSocket: {
        port: 80
      },
      exec: {
        command: ['']
      },
      initialDelaySeconds: 0,
      periodSeconds: 10,
      timeoutSeconds: 1,
      successThreshold: 1,
      failureThreshold: 3
    },
    livenessProbe: {
      type: 'httpGet',
      httpGet: {
        path: '',
        port: 80,
        scheme: 'HTTP'
      },
      tcpSocket: {
        port: 80
      },
      exec: {
        command: ['']
      },
      initialDelaySeconds: 0,
      periodSeconds: 10,
      timeoutSeconds: 1,
      successThreshold: 1,
      failureThreshold: 3
    },
    readinessProbe: {
      type: 'httpGet',
      httpGet: {
        path: '',
        port: 80,
        scheme: 'HTTP'
      },
      tcpSocket: {
        port: 80
      },
      exec: {
        command: ['']
      },
      initialDelaySeconds: 0,
      periodSeconds: 10,
      timeoutSeconds: 1,
      successThreshold: 1,
      failureThreshold: 3
    }
  },
  // 环境变量
  env: [{
    name: '',
    value: '',
    valueFrom: {
      configMapKeyRef: {
        name: '',
        key: ''
      },
      secretKeyRef: {
        name: '',
        key: ''
      }
    },
    type: 'value'
  }],
  // 数据存储
  volumes: [{
    name: '',
    type: 'emptyDir',
    emptyDir: {
      medium: ''
    },
    hostPath: {
      path: '',
      type: ''
    },
    configMap: {
      name: '',
      items: [{
        key: '',
        path: ''
      }]
    },
    secret: {
      secretName: '',
      items: [{
        key: '',
        path: ''
      }]
    },
    persistentVolumeClaim: {
      claimName: ''
    }
  }],
  volumeMounts: [{
    name: '',
    mountPath: '',
    subPath: '',
    readOnly: false
  }],
  // Init容器
  initContainers: [{
    name: '',
    image: '',
    command: [''],
    args: [''],
    env: [{
      name: '',
      value: '',
      valueFrom: {
        configMapKeyRef: {
          name: '',
          key: ''
        },
        secretKeyRef: {
          name: '',
          key: ''
        }
      },
      type: 'value'
    }],
    resources: {
      requests: {
        cpu: '100m',
        memory: '128Mi'
      },
      limits: {
        cpu: '500m',
        memory: '256Mi'
      }
    },
    volumeMounts: [{
      name: '',
      mountPath: '',
      subPath: '',
      readOnly: false
    }]
  }],
  // 镜像拉取凭证
  imagePullSecrets: [{
    name: ''
  }],
  // 升级策略
  strategy: {
    type: 'RollingUpdate',
    rollingUpdate: {
      maxSurge: '25%',
      maxUnavailable: '25%'
    }
  },
  // 调度策略
  nodeSelector: [{
    key: '',
    value: ''
  }],
  affinity: {
    nodeAffinity: {
      requiredDuringSchedulingIgnoredDuringExecution: {
        nodeSelectorTerms: [{
          matchExpressions: [{
            key: '',
            operator: 'In',
            values: ['']
          }]
        }]
      },
      preferredDuringSchedulingIgnoredDuringExecution: [{
        weight: 1,
        preference: {
          matchExpressions: [{
            key: '',
            operator: 'In',
            values: ['']
          }]
        }
      }]
    },
    podAffinity: {
      requiredDuringSchedulingIgnoredDuringExecution: [{
        labelSelector: {
          matchExpressions: [{
            key: '',
            operator: 'In',
            values: ['']
          }]
        },
        topologyKey: ''
      }],
      preferredDuringSchedulingIgnoredDuringExecution: [{
        weight: 1,
        podAffinityTerm: {
          labelSelector: {
            matchExpressions: [{
              key: '',
              operator: 'In',
              values: ['']
            }]
          },
          topologyKey: ''
        }
      }]
    },
    podAntiAffinity: {
      requiredDuringSchedulingIgnoredDuringExecution: [{
        labelSelector: {
          matchExpressions: [{
            key: '',
            operator: 'In',
            values: ['']
          }]
        },
        topologyKey: ''
      }],
      preferredDuringSchedulingIgnoredDuringExecution: [{
        weight: 1,
        podAffinityTerm: {
          labelSelector: {
            matchExpressions: [{
              key: '',
              operator: 'In',
              values: ['']
            }]
          },
          topologyKey: ''
        }
      }]
    }
  },
  // 容忍策略
  tolerations: [{
    key: '',
    operator: 'Equal',
    value: '',
    effect: 'NoSchedule'
  }],
  // 标签
  labels: [{
    key: '',
    value: ''
  }],
  // 注解
  annotations: [{
    key: '',
    value: ''
  }],
  // 资源配置
  resources: {
    requests: {
      cpu: '100m',
      memory: '128Mi'
    },
    limits: {
      cpu: '500m',
      memory: '256Mi'
    }
  },
  // DNS配置
  dnsPolicy: 'ClusterFirst',
  dnsConfig: {
    nameservers: [''],
    searches: [''],
    options: [{
      name: '',
      value: ''
    }]
  },
  // 其他配置
  terminationGracePeriodSeconds: 30,
  hostNetwork: false
})



// 组件挂载时加载数据
onMounted(() => {
  loadClusters()
})

// 监听集群选择变化，获取对应的命名空间列表
watch(selectedCluster, (newClusterId) => {
  if (newClusterId) {
    loadNamespaces(newClusterId)
  } else {
    // 重置命名空间列表
    namespaces.value = []
    selectedNamespace.value = ''
  }
})

// 监听命名空间选择变化，自动刷新当前工作负载列表
watch(selectedNamespace, () => {
  if (selectedCluster.value) {
    loadCurrentWorkload()
  }
})

// 加载集群列表
const loadClusters = async () => {
  try {
    // 调用真实的API获取集群列表，axios会自动添加baseURL
    const response = await axios.get('/clusters')
    clusters.value = response.data
  } catch (error) {
    console.error('加载集群列表失败:', error)
    // 不显示错误消息给用户，只在控制台打印
    clusters.value = []
  }
}

// 加载命名空间列表
const loadNamespaces = async (clusterId: number) => {
  try {
    namespacesLoading.value = true
    // 调用API获取指定集群的命名空间列表，axios会自动添加baseURL
    // 添加分页参数，获取所有命名空间
    const response = await axios.get(`/clusters/${clusterId}/namespaces`, {
      params: {
        page: 1,
        pageSize: 1000 // 设置一个足够大的值，确保获取所有命名空间
      }
    })
    // 检查后端返回的数据格式，如果是分页对象则使用data字段
    const namespacesData = response.data.data || response.data
    // 只提取命名空间名称，确保是字符串数组
    namespaces.value = namespacesData.map((ns: any) => ns.name || '')
    // 过滤掉空字符串
    namespaces.value = namespaces.value.filter((ns: string) => ns !== '')
    // 不自动选择命名空间，让用户手动选择
    selectedNamespace.value = ''
    namespacesLoading.value = false
  } catch (error) {
    console.error('加载命名空间列表失败:', error)
    // 不显示错误消息给用户，只在控制台打印
    namespaces.value = []
    selectedNamespace.value = ''
    namespacesLoading.value = false
  }
}

// 取消令牌，用于取消之前的请求
let cancelTokenSource: any = null

// 加载Deployment列表
const loadDeployments = async () => {
  if (!selectedCluster.value) {
    ElMessage.warning('请先选择一个集群')
    return
  }
  
  // 取消之前的请求
  if (cancelTokenSource) {
    cancelTokenSource.cancel('Operation canceled by the user.')
  }
  
  // 创建新的取消令牌
  cancelTokenSource = axios.CancelToken.source()
  
  loading.value.deployment = true
  try {
    // 调用API获取Deployment列表，包含命名空间参数，axios会自动添加baseURL
    const response = await axios.get(`/clusters/${selectedCluster.value}/deployments`, {
      params: {
        namespace: selectedNamespace.value || '',
        _t: Date.now() // 添加时间戳，防止浏览器缓存
      },
      headers: {
        'Cache-Control': 'no-cache' // 禁用缓存
      },
      cancelToken: cancelTokenSource.token
    })
    // 为每个Deployment添加集群信息
    const clusterName = clusters.value.find(c => c.id === selectedCluster.value)?.name || ''
    // 查看API返回的原始createdAt格式
    if (response.data.length > 0) {
      console.log('API返回的createdAt原始格式:', response.data[0].createdAt, typeof response.data[0].createdAt)
    }
    const newDeployments = response.data.map((deployment: any) => ({
      ...deployment,
      cluster: clusterName,
      clusterID: selectedCluster.value
    }))
    
    // 先更新数据，然后设置loading为false，确保数据正确显示
    deployments.value = newDeployments
    loading.value.deployment = false
    
    // 如果返回的Deployment列表为空，显示提示信息
    if (newDeployments.length === 0) {
      ElMessage.info('目标命名空间内无Deployment实例')
    }
  } catch (error) {
    // 如果是取消请求的错误，不显示错误信息
    if (axios.isCancel(error)) {
      console.log('Request canceled:', error.message)
      return
    }
    console.error('加载Deployment列表失败:', error)
    loading.value.deployment = false
    ElMessage.error('加载Deployment列表失败')
    // 确保在错误情况下，列表为空，显示正确的空状态
    deployments.value = []
  }
}

// 加载StatefulSet列表
const loadStatefulSets = async () => {
  if (!selectedCluster.value) {
    ElMessage.warning('请先选择一个集群')
    return
  }
  
  // 取消之前的请求
  if (cancelTokenSource) {
    cancelTokenSource.cancel('Operation canceled by the user.')
  }
  
  // 创建新的取消令牌
  cancelTokenSource = axios.CancelToken.source()
  
  loading.value.statefulset = true
  try {
    // 调用API获取StatefulSet列表，包含命名空间参数，axios会自动添加baseURL
    const response = await axios.get(`/clusters/${selectedCluster.value}/statefulsets`, {
      params: {
        namespace: selectedNamespace.value || '',
        _t: Date.now() // 添加时间戳，防止浏览器缓存
      },
      headers: {
        'Cache-Control': 'no-cache' // 禁用缓存
      },
      cancelToken: cancelTokenSource.token
    })
    // 为每个StatefulSet添加集群信息
    const clusterName = clusters.value.find(c => c.id === selectedCluster.value)?.name || ''
    const newStatefulSets = response.data.map((statefulset: any) => ({
      ...statefulset,
      cluster: clusterName,
      clusterID: selectedCluster.value
    }))
    
    // 先更新数据，然后设置loading为false，确保数据正确显示
    statefulsets.value = newStatefulSets
    loading.value.statefulset = false
    
    // 如果返回的StatefulSet列表为空，显示提示信息
    if (newStatefulSets.length === 0) {
      ElMessage.info('目标命名空间内无StatefulSet实例')
    }
  } catch (error) {
    // 如果是取消请求的错误，不显示错误信息
    if (axios.isCancel(error)) {
      console.log('Request canceled:', error.message)
      return
    }
    console.error('加载StatefulSet列表失败:', error)
    loading.value.statefulset = false
    ElMessage.error('加载StatefulSet列表失败')
    // 确保在错误情况下，列表为空，显示正确的空状态
    statefulsets.value = []
  }
}

// 加载DaemonSet列表
const loadDaemonSets = async () => {
  if (!selectedCluster.value) {
    ElMessage.warning('请先选择一个集群')
    return
  }
  
  // 取消之前的请求
  if (cancelTokenSource) {
    cancelTokenSource.cancel('Operation canceled by the user.')
  }
  
  // 创建新的取消令牌
  cancelTokenSource = axios.CancelToken.source()
  
  loading.value.daemonset = true
  try {
    // 调用API获取DaemonSet列表，包含命名空间参数，axios会自动添加baseURL
    const response = await axios.get(`/clusters/${selectedCluster.value}/daemonsets`, {
      params: {
        namespace: selectedNamespace.value || '',
        _t: Date.now() // 添加时间戳，防止浏览器缓存
      },
      headers: {
        'Cache-Control': 'no-cache' // 禁用缓存
      },
      cancelToken: cancelTokenSource.token
    })
    // 为每个DaemonSet添加集群信息
    const clusterName = clusters.value.find(c => c.id === selectedCluster.value)?.name || ''
    const newDaemonSets = response.data.map((daemonset: any) => ({
      ...daemonset,
      cluster: clusterName,
      clusterID: selectedCluster.value
    }))
    
    // 先更新数据，然后设置loading为false，确保数据正确显示
    daemonsets.value = newDaemonSets
    loading.value.daemonset = false
    
    // 如果返回的DaemonSet列表为空，显示提示信息
    if (newDaemonSets.length === 0) {
      ElMessage.info('目标命名空间内无DaemonSet实例')
    }
  } catch (error) {
    // 如果是取消请求的错误，不显示错误信息
    if (axios.isCancel(error)) {
      console.log('Request canceled:', error.message)
      return
    }
    console.error('加载DaemonSet列表失败:', error)
    loading.value.daemonset = false
    ElMessage.error('加载DaemonSet列表失败')
    // 确保在错误情况下，列表为空，显示正确的空状态
    daemonsets.value = []
  }
}

// 加载Job列表
const loadJobs = async () => {
  if (!selectedCluster.value) {
    ElMessage.warning('请先选择一个集群')
    return
  }
  
  // 取消之前的请求
  if (cancelTokenSource) {
    cancelTokenSource.cancel('Operation canceled by the user.')
  }
  
  // 创建新的取消令牌
  cancelTokenSource = axios.CancelToken.source()
  
  loading.value.job = true
  try {
    // 调用API获取Job列表，包含命名空间参数，axios会自动添加baseURL
    const response = await axios.get(`/clusters/${selectedCluster.value}/jobs`, {
      params: {
        namespace: selectedNamespace.value || '',
        _t: Date.now() // 添加时间戳，防止浏览器缓存
      },
      headers: {
        'Cache-Control': 'no-cache' // 禁用缓存
      },
      cancelToken: cancelTokenSource.token
    })
    // 为每个Job添加集群信息
    const clusterName = clusters.value.find(c => c.id === selectedCluster.value)?.name || ''
    const newJobs = response.data.map((job: any) => ({
      ...job,
      cluster: clusterName,
      clusterID: selectedCluster.value
    }))
    
    // 先更新数据，然后设置loading为false，确保数据正确显示
    jobs.value = newJobs
    loading.value.job = false
    
    // 如果返回的Job列表为空，显示提示信息
    if (newJobs.length === 0) {
      ElMessage.info('目标命名空间内无Job实例')
    }
  } catch (error) {
    // 如果是取消请求的错误，不显示错误信息
    if (axios.isCancel(error)) {
      console.log('Request canceled:', error.message)
      return
    }
    console.error('加载Job列表失败:', error)
    loading.value.job = false
    ElMessage.error('加载Job列表失败')
    // 确保在错误情况下，列表为空，显示正确的空状态
    jobs.value = []
  }
}

// 加载CronJob列表
const loadCronJobs = async () => {
  if (!selectedCluster.value) {
    ElMessage.warning('请先选择一个集群')
    return
  }
  
  // 取消之前的请求
  if (cancelTokenSource) {
    cancelTokenSource.cancel('Operation canceled by the user.')
  }
  
  // 创建新的取消令牌
  cancelTokenSource = axios.CancelToken.source()
  
  loading.value.cronjob = true
  try {
    // 调用API获取CronJob列表，包含命名空间参数，axios会自动添加baseURL
    const response = await axios.get(`/clusters/${selectedCluster.value}/cronjobs`, {
      params: {
        namespace: selectedNamespace.value || '',
        _t: Date.now() // 添加时间戳，防止浏览器缓存
      },
      headers: {
        'Cache-Control': 'no-cache' // 禁用缓存
      },
      cancelToken: cancelTokenSource.token
    })
    // 为每个CronJob添加集群信息
    const clusterName = clusters.value.find(c => c.id === selectedCluster.value)?.name || ''
    const newCronJobs = response.data.map((cronjob: any) => ({
      ...cronjob,
      cluster: clusterName,
      clusterID: selectedCluster.value
    }))
    
    // 先更新数据，然后设置loading为false，确保数据正确显示
    cronjobs.value = newCronJobs
    loading.value.cronjob = false
    
    // 如果返回的CronJob列表为空，显示提示信息
    if (newCronJobs.length === 0) {
      ElMessage.info('目标命名空间内无CronJob实例')
    }
  } catch (error) {
    // 如果是取消请求的错误，不显示错误信息
    if (axios.isCancel(error)) {
      console.log('Request canceled:', error.message)
      return
    }
    console.error('加载CronJob列表失败:', error)
    loading.value.cronjob = false
    ElMessage.error('加载CronJob列表失败')
    // 确保在错误情况下，列表为空，显示正确的空状态
    cronjobs.value = []
  }
}

// 根据当前标签页加载对应的工作负载数据
const loadCurrentWorkload = () => {
  switch (activeWorkloadTab.value) {
    case 'deployment':
      loadDeployments()
      break
    case 'statefulset':
      loadStatefulSets()
      break
    case 'daemonset':
      loadDaemonSets()
      break
    case 'job':
      loadJobs()
      break
    case 'cronjob':
      loadCronJobs()
      break
    default:
      break
  }
}

// 刷新当前工作负载数据
const refreshCurrentWorkload = () => {
  loadCurrentWorkload()
}

// 标签页切换处理
const handleTabChange = () => {
  loadCurrentWorkload()
}

// 获取状态标签类型
const getStatusTagType = (status: string) => {
  switch (status) {
    case 'Available':
    case 'Ready':
    case 'Running':
      return 'success'
    case 'Progressing':
    case 'Pending':
      return 'warning'
    case 'Degraded':
    case 'Failed':
      return 'danger'
    case 'Unknown':
      return 'info'
    default:
      return 'info'
  }
}

// 格式化绝对时间
// 格式化完整日期时间
const formatDateTime = (datetime: string | number) => {
  // 空值处理
  if (!datetime) {
    return 'No Data'
  }
  
  try {
    let date: Date
    
    // 处理不同类型的时间输入
    if (typeof datetime === 'string') {
      // 处理ISO格式时间（包含T和Z）
      if (datetime.includes('T') || datetime.includes('Z')) {
        date = new Date(datetime)
      } 
      // 处理其他字符串格式
      else {
        // 尝试直接解析
        date = new Date(datetime)
        
        // 如果解析失败，尝试转换为时间戳
        if (isNaN(date.getTime())) {
          const timestamp = parseInt(datetime)
          if (!isNaN(timestamp)) {
            date = new Date(timestamp)
          }
        }
      }
    } else if (typeof datetime === 'number') {
      // 处理时间戳（秒或毫秒）
      date = new Date(datetime)
      
      // 检查是否是秒级时间戳（小于13位是秒，大于等于13位是毫秒）
      if (datetime.toString().length < 13) {
        date = new Date(datetime * 1000)
      }
    } else {
      // 处理其他类型
      date = new Date(datetime.toString())
    }
    
    // 检查日期是否有效
    if (isNaN(date.getTime())) {
      return 'Invalid Date'
    }
    
    // 格式化日期为YYYY-MM-DD HH:MM:SS格式
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hours = String(date.getHours()).padStart(2, '0')
    const minutes = String(date.getMinutes()).padStart(2, '0')
    const seconds = String(date.getSeconds()).padStart(2, '0')
    
    // 检查是否是默认的无效日期（1970-01-01或0001-01-01）
    if ((year === 1970 && month === '01' && day === '01') || 
        (year === 1 && month === '01' && day === '01')) {
      return 'Invalid Date'
    }
    
    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
  } catch (error) {
    return 'Invalid Date'
  }
}

// 保留原有的formatAge函数，用于其他地方
const formatAge = (createdAt: string) => {
  return formatDateTime(createdAt)
}

// 查看工作负载详情
const viewWorkloadDetails = (type: string, _workload: any) => {
  // 这里可以根据不同的工作负载类型显示不同的详情对话框
  ElMessage.info(`查看${type}详情功能待实现`)
}

// 处理工作负载命令
const handleWorkloadCommand = (command: string, type: string, workload: any) => {
  switch (command) {
    case 'scale':
      if (type === 'deployment' || type === 'statefulset') {
        scaleWorkload(type, workload)
      } else {
        ElMessage.info(`扩缩容功能仅支持Deployment和StatefulSet`)
      }
      break
    case 'restart':
      if (type === 'deployment' || type === 'statefulset' || type === 'daemonset') {
        restartWorkload(type, workload)
      } else {
        ElMessage.info(`重启功能仅支持Deployment、StatefulSet和DaemonSet`)
      }
      break
    case 'versions':
      if (type === 'deployment') {
        viewReplicaSetVersions(workload)
      } else {
        ElMessage.info(`版本记录功能仅支持Deployment`)
      }
      break
    case 'delete':
      deleteWorkload(type, workload)
      break
    case 'run':
      if (type === 'cronjob') {
        runCronJob(workload)
      } else {
        ElMessage.info(`立即运行功能仅支持CronJob`)
      }
      break
    default:
      break
  }
}

// 查看ReplicaSet版本记录
const viewReplicaSetVersions = (workload: any) => {
  selectedDeployment.value = workload
  loadReplicaSets(workload)
  replicasetVersionsVisible.value = true
}

// 加载ReplicaSet版本记录
const loadReplicaSets = async (workload: any) => {
  if (!selectedCluster.value) {
    ElMessage.warning('请先选择一个集群')
    return
  }

  try {
    // 调用API获取ReplicaSet列表，axios会自动添加baseURL
    const response = await axios.get(`/clusters/${selectedCluster.value}/replicasets`, {
      params: {
        namespace: workload.namespace,
        deployment: workload.name,
        _t: Date.now() // 添加时间戳，防止浏览器缓存
      },
      headers: {
        'Cache-Control': 'no-cache' // 禁用缓存
      }
    })
    
    // 为每个ReplicaSet添加基本信息
    const newReplicaSets = response.data.map((rs: any) => ({
      name: rs.Name,
      namespace: rs.Namespace,
      replicas: rs.Replicas,
      readyReplicas: rs.ReadyReplicas,
      image: rs.Image,
      createdAt: rs.CreatedAt,
      deployment: workload.name
    }))
    
    // 按创建时间倒序排序，最新的版本在前
    newReplicaSets.sort((a: ReplicaSet, b: ReplicaSet) => {
      return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
    })
    
    replicasets.value = newReplicaSets
  } catch (error) {
    console.error('加载ReplicaSet列表失败:', error)
    ElMessage.error('加载ReplicaSet列表失败')
    replicasets.value = []
  }
}

// 扩缩容工作负载
const scaleWorkload = (_type: string, workload: any) => {
  // 设置选中的工作负载和初始副本数
  selectedDeployment.value = workload
  scaleForm.value.replicas = workload.replicas
  // 打开扩缩容对话框
  scaleDeploymentVisible.value = true
}

// 处理扩缩容操作
const handleScaleDeployment = async () => {
  if (!selectedDeployment.value) return
  
  try {
    // 调用扩缩容API
    console.log('调用扩缩容API:', {
      url: `/clusters/${selectedDeployment.value.clusterID}/deployments/${selectedDeployment.value.name}/scale`,
      data: {
        namespace: selectedDeployment.value.namespace,
        replicas: scaleForm.value.replicas
      }
    })
    
    const response = await axios.put(`/clusters/${selectedDeployment.value.clusterID}/deployments/${selectedDeployment.value.name}/scale`, {
      namespace: selectedDeployment.value.namespace,
      replicas: scaleForm.value.replicas
    })
    
    console.log('扩缩容API响应:', response.data)
    
    // 关闭对话框
    scaleDeploymentVisible.value = false
    
    // 刷新当前工作负载列表
    loadCurrentWorkload()
    
    // 显示成功消息
    ElMessage.success('扩缩容成功')
  } catch (error: any) {
    console.error('扩缩容失败:', error)
    console.error('错误详情:', {
      message: error.message,
      response: error.response?.data,
      status: error.response?.status,
      headers: error.response?.headers
    })
    ElMessage.error(`扩缩容失败: ${error.response?.data?.error || error.message || '未知错误'}`)
  }
}

// 重启工作负载
const restartWorkload = async (type: string, workload: any) => {
  try {
    // 目前只支持Deployment重启
    if (type === 'deployment') {
      // 调用重启API
      await axios.post(`/clusters/${selectedCluster.value}/deployments/${workload.name}/restart`, {
        namespace: workload.namespace
      })
      
      // 刷新当前工作负载列表
      loadCurrentWorkload()
      
      // 显示成功消息
      ElMessage.success(`${type}重启成功`)
    } else {
      // 其他类型暂不支持
      ElMessage.info(`${type}重启功能开发中`)
    }
  } catch (error) {
    console.error(`${type}重启失败:', error`)
    ElMessage.error(`${type}重启失败`)
  }
}

// 删除工作负载
const deleteWorkload = async (type: string, workload: any) => {
  try {
    // 这里应该调用真实的API删除工作负载
    // 暂时使用模拟数据
    let index = -1
    switch (type) {
      case 'deployment':
        index = deployments.value.findIndex((d: any) => d.namespace === workload.namespace && d.name === workload.name)
        if (index !== -1) {
          deployments.value.splice(index, 1)
        }
        break
      case 'statefulset':
        index = statefulsets.value.findIndex((s: any) => s.namespace === workload.namespace && s.name === workload.name)
        if (index !== -1) {
          statefulsets.value.splice(index, 1)
        }
        break
      case 'daemonset':
        index = daemonsets.value.findIndex((d: any) => d.namespace === workload.namespace && d.name === workload.name)
        if (index !== -1) {
          daemonsets.value.splice(index, 1)
        }
        break
      case 'job':
        index = jobs.value.findIndex((j: any) => j.namespace === workload.namespace && j.name === workload.name)
        if (index !== -1) {
          jobs.value.splice(index, 1)
        }
        break
      case 'cronjob':
        index = cronjobs.value.findIndex((c: any) => c.namespace === workload.namespace && c.name === workload.name)
        if (index !== -1) {
          cronjobs.value.splice(index, 1)
        }
        break
    }
    ElMessage.success(`删除${type}成功`)
  } catch (error) {
    console.error(`删除${type}失败:', error`)
    ElMessage.error(`删除${type}失败`)
  }
}

// 立即运行CronJob
const runCronJob = async (_cronjob: any) => {
  try {
    // 这里应该调用真实的API立即运行CronJob
    // 暂时使用模拟数据
    ElMessage.success('立即运行CronJob成功')
    // 运行后刷新列表
    loadCronJobs()
  } catch (error) {
    console.error('立即运行CronJob失败:', error)
    ElMessage.error('立即运行CronJob失败')
  }
}

// 创建工作负载
const createWorkload = () => {
  if (activeWorkloadTab.value === 'deployment') {
    // 重置表单和YAML内容
    resetCreateForm()
    createDeploymentVisible.value = true
  } else {
    ElMessage.info(`创建${activeWorkloadTab.value}功能待实现`)
  }
}

// 重置创建表单
const resetCreateForm = () => {
  createMode.value = 'form'
  yamlContent.value = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: example
  namespace: default
  labels:
    app: app
spec:
  replicas: 1
  selector:
    matchLabels:
      app: app
  template:
    metadata:
      labels:
        app: app
    spec:
      containers:
        - name: main
          image: nginx:latest
          imagePullPolicy: IfNotPresent
          resources:
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 500m
              memory: 512Mi`
  validationResult.value = null
  deploymentForm.value = {
    name: '',
    namespace: selectedNamespace.value || '',
    replicas: 1,
    image: '',
    containerName: 'app',
    ports: [{
      containerPort: 80,
      protocol: 'TCP'
    }],
    // 生命周期配置
    lifecycle: {
      startupProbe: {
        type: 'httpGet',
        httpGet: {
          path: '',
          port: 80,
          scheme: 'HTTP'
        },
        tcpSocket: {
          port: 80
        },
        exec: {
          command: ['']
        },
        initialDelaySeconds: 0,
        periodSeconds: 10,
        timeoutSeconds: 1,
        successThreshold: 1,
        failureThreshold: 3
      },
      livenessProbe: {
        type: 'httpGet',
        httpGet: {
          path: '',
          port: 80,
          scheme: 'HTTP'
        },
        tcpSocket: {
          port: 80
        },
        exec: {
          command: ['']
        },
        initialDelaySeconds: 0,
        periodSeconds: 10,
        timeoutSeconds: 1,
        successThreshold: 1,
        failureThreshold: 3
      },
      readinessProbe: {
        type: 'httpGet',
        httpGet: {
          path: '',
          port: 80,
          scheme: 'HTTP'
        },
        tcpSocket: {
          port: 80
        },
        exec: {
          command: ['']
        },
        initialDelaySeconds: 0,
        periodSeconds: 10,
        timeoutSeconds: 1,
        successThreshold: 1,
        failureThreshold: 3
      }
    },
    // 环境变量
    env: [{
      name: '',
      value: '',
      valueFrom: {
        configMapKeyRef: {
          name: '',
          key: ''
        },
        secretKeyRef: {
          name: '',
          key: ''
        }
      },
      type: 'value'
    }],
    // 数据存储
    volumes: [{
      name: '',
      type: 'emptyDir',
      emptyDir: {
        medium: ''
      },
      hostPath: {
        path: '',
        type: ''
      },
      configMap: {
        name: '',
        items: [{
          key: '',
          path: ''
        }]
      },
      secret: {
        secretName: '',
        items: [{
          key: '',
          path: ''
        }]
      },
      persistentVolumeClaim: {
        claimName: ''
      }
    }],
    volumeMounts: [{
      name: '',
      mountPath: '',
      subPath: '',
      readOnly: false
    }],
    // Init容器
    initContainers: [{
      name: '',
      image: '',
      command: [''],
      args: [''],
      env: [{
        name: '',
        value: '',
        valueFrom: {
          configMapKeyRef: {
            name: '',
            key: ''
          },
          secretKeyRef: {
            name: '',
            key: ''
          }
        },
        type: 'value'
      }],
      resources: {
        requests: {
          cpu: '100m',
          memory: '128Mi'
        },
        limits: {
          cpu: '500m',
          memory: '256Mi'
        }
      },
      volumeMounts: [{
        name: '',
        mountPath: '',
        subPath: '',
        readOnly: false
      }]
    }],
    // 镜像拉取凭证
    imagePullSecrets: [{
      name: ''
    }],
    // 升级策略
    strategy: {
      type: 'RollingUpdate',
      rollingUpdate: {
        maxSurge: '25%',
        maxUnavailable: '25%'
      }
    },
    // 调度策略
    nodeSelector: [{
      key: '',
      value: ''
    }],
    affinity: {
      nodeAffinity: {
        requiredDuringSchedulingIgnoredDuringExecution: {
          nodeSelectorTerms: [{
            matchExpressions: [{
              key: '',
              operator: 'In',
              values: ['']
            }]
          }]
        },
        preferredDuringSchedulingIgnoredDuringExecution: [{
          weight: 1,
          preference: {
            matchExpressions: [{
              key: '',
              operator: 'In',
              values: ['']
            }]
          }
        }]
      },
      podAffinity: {
        requiredDuringSchedulingIgnoredDuringExecution: [{
          labelSelector: {
            matchExpressions: [{
              key: '',
              operator: 'In',
              values: ['']
            }]
          },
          topologyKey: ''
        }],
        preferredDuringSchedulingIgnoredDuringExecution: [{
          weight: 1,
          podAffinityTerm: {
            labelSelector: {
              matchExpressions: [{
                key: '',
                operator: 'In',
                values: ['']
              }]
            },
            topologyKey: ''
          }
        }]
      },
      podAntiAffinity: {
        requiredDuringSchedulingIgnoredDuringExecution: [{
          labelSelector: {
            matchExpressions: [{
              key: '',
              operator: 'In',
              values: ['']
            }]
          },
          topologyKey: ''
        }],
        preferredDuringSchedulingIgnoredDuringExecution: [{
          weight: 1,
          podAffinityTerm: {
            labelSelector: {
              matchExpressions: [{
                key: '',
                operator: 'In',
                values: ['']
              }]
            },
            topologyKey: ''
          }
        }]
      }
    },
    // 容忍策略
    tolerations: [{
      key: '',
      operator: 'Equal',
      value: '',
      effect: 'NoSchedule'
    }],
    // 标签
    labels: [{
      key: '',
      value: ''
    }],
    // 注解
    annotations: [{
      key: '',
      value: ''
    }],
    // 资源配置
    resources: {
      requests: {
        cpu: '100m',
        memory: '128Mi'
      },
      limits: {
        cpu: '500m',
        memory: '256Mi'
      }
    },
    // DNS配置
    dnsPolicy: 'ClusterFirst',
    dnsConfig: {
      nameservers: [''],
      searches: [''],
      options: [{
        name: '',
        value: ''
      }]
    },
    // 其他配置
    terminationGracePeriodSeconds: 30,
    hostNetwork: false
  }
}

// 切换创建模式
const toggleCreateMode = () => {
  // 切换模式时清空验证结果
  validationResult.value = null
  // 只有从表单模式切换到YAML模式，且YAML内容为空时，才从表单生成YAML
  // 否则保留现有YAML内容
  if (createMode.value === 'yaml' && yamlContent.value === '') {
    // 从表单生成YAML
    generateYamlFromForm()
  }
}

// 从表单生成YAML
const generateYamlFromForm = () => {
  const form = deploymentForm.value
  let yaml = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: ${form.name}
  namespace: ${form.namespace}
  labels:
${form.labels.filter(l => l.key && l.value).map(l => `    ${l.key}: ${l.value}`).join('\n')}
  annotations:
${form.annotations.filter(a => a.key && a.value).map(a => `    ${a.key}: ${a.value}`).join('\n')}
spec:
  replicas: ${form.replicas}
  strategy:
    type: ${form.strategy.type}`

  // 添加滚动更新配置
  if (form.strategy.type === 'RollingUpdate') {
    yaml += `
    rollingUpdate:
      maxUnavailable: ${form.strategy.rollingUpdate.maxUnavailable}
      maxSurge: ${form.strategy.rollingUpdate.maxSurge}`
  }

  yaml += `
  selector:
    matchLabels:
${form.labels.filter(l => l.key && l.value).map(l => `      ${l.key}: ${l.value}`).join('\n')}
  template:
    metadata:
      labels:
${form.labels.filter(l => l.key && l.value).map(l => `        ${l.key}: ${l.value}`).join('\n')}
    spec:
      dnsPolicy: ${form.dnsPolicy}`

  // 添加DNS配置
  if (form.dnsConfig.nameservers[0] || form.dnsConfig.searches[0] || form.dnsConfig.options.length > 0) {
    yaml += `
      dnsConfig:`
    if (form.dnsConfig.nameservers[0]) {
      yaml += `
        nameservers:
${form.dnsConfig.nameservers[0].split(',').map(ns => `        - ${ns.trim()}`).join('\n')}`
    }
    if (form.dnsConfig.searches[0]) {
      yaml += `
        searches:
${form.dnsConfig.searches[0].split(',').map(search => `        - ${search.trim()}`).join('\n')}`
    }
    if (form.dnsConfig.options.length > 0) {
      yaml += `
        options:`
      for (const option of form.dnsConfig.options) {
        if (option.name) {
          yaml += `
        - name: ${option.name}`
          if (option.value) {
            yaml += `
          value: ${option.value}`
          }
        }
      }
    }
  }

  // 添加镜像拉取凭证
  if (form.imagePullSecrets.length > 0) {
    yaml += `
      imagePullSecrets:
${form.imagePullSecrets.map(secret => `      - name: ${secret.name}`).join('\n')}`
  }

  // 添加主机网络
  if (form.hostNetwork) {
    yaml += `
      hostNetwork: true`
  }

  // 添加Node Selector
  const nodeSelectors = form.nodeSelector.filter(s => s.key && s.value)
  if (nodeSelectors.length > 0) {
    yaml += `
      nodeSelector:
${nodeSelectors.map(s => `        ${s.key}: ${s.value}`).join('\n')}`
  }

  // 添加容忍策略
  if (form.tolerations.length > 0) {
    yaml += `
      tolerations:`
    for (const toleration of form.tolerations) {
      if (toleration.key) {
        yaml += `
      - key: ${toleration.key}`
        yaml += `
        operator: ${toleration.operator}`
        if (toleration.value) {
          yaml += `
        value: ${toleration.value}`
        }
        yaml += `
        effect: ${toleration.effect}`
      }
    }
  }

  // 添加卷配置
  const validVolumes = form.volumes.filter(v => v.name && v.type)
  if (validVolumes.length > 0) {
    yaml += `
      volumes:`
    for (const volume of validVolumes) {
      yaml += `
      - name: ${volume.name}`
      if (volume.type === 'emptyDir') {
        yaml += `
        emptyDir:`
        if (volume.emptyDir.medium) {
          yaml += `
          medium: ${volume.emptyDir.medium}`
        }
      } else if (volume.type === 'hostPath') {
        yaml += `
        hostPath:`
        yaml += `
          path: ${volume.hostPath.path}`
        if (volume.hostPath.type) {
          yaml += `
          type: ${volume.hostPath.type}`
        }
      } else if (volume.type === 'configMap') {
        yaml += `
        configMap:`
        yaml += `
          name: ${volume.configMap.name}`
      } else if (volume.type === 'secret') {
        yaml += `
        secret:`
        yaml += `
          secretName: ${volume.secret.secretName}`
      } else if (volume.type === 'persistentVolumeClaim') {
        yaml += `
        persistentVolumeClaim:`
        yaml += `
          claimName: ${volume.persistentVolumeClaim.claimName}`
      }
    }
  }

  // 添加Init容器
  if (form.initContainers.length > 0) {
    for (const container of form.initContainers) {
      if (container.name && container.image) {
        yaml += `
      initContainers:
      - name: ${container.name}
        image: ${container.image}`
        if (container.command[0]) {
          yaml += `
        command:
        - ${container.command[0]}`
        }
        if (container.args[0]) {
          yaml += `
        args:
        - ${container.args[0]}`
        }
        // 添加Init容器环境变量
        const initEnvs = container.env.filter(env => env.name && (env.type === 'value' ? env.value : env.type === 'configMapKeyRef' ? env.valueFrom.configMapKeyRef.name && env.valueFrom.configMapKeyRef.key : env.valueFrom.secretKeyRef.name && env.valueFrom.secretKeyRef.key))
        if (initEnvs.length > 0) {
          yaml += `
        env:`
          for (const env of initEnvs) {
            yaml += `
        - name: ${env.name}`
            if (env.type === 'value') {
              yaml += `
          value: ${env.value}`
            } else if (env.type === 'configMapKeyRef') {
              yaml += `
          valueFrom:
            configMapKeyRef:
              name: ${env.valueFrom.configMapKeyRef.name}
              key: ${env.valueFrom.configMapKeyRef.key}`
            } else if (env.type === 'secretKeyRef') {
              yaml += `
          valueFrom:
            secretKeyRef:
              name: ${env.valueFrom.secretKeyRef.name}
              key: ${env.valueFrom.secretKeyRef.key}`
            }
          }
        }
        // 添加Init容器资源配置
        yaml += `
        resources:
          requests:
            cpu: ${container.resources.requests.cpu}
            memory: ${container.resources.requests.memory}
          limits:
            cpu: ${container.resources.limits.cpu}
            memory: ${container.resources.limits.memory}`
        // 添加Init容器卷挂载
        const initVolumeMounts = container.volumeMounts.filter(mount => mount.name && mount.mountPath)
        if (initVolumeMounts.length > 0) {
          yaml += `
        volumeMounts:`
          for (const mount of initVolumeMounts) {
            yaml += `
        - name: ${mount.name}`
            yaml += `
          mountPath: ${mount.mountPath}`
            if (mount.subPath) {
              yaml += `
          subPath: ${mount.subPath}`
            }
            if (mount.readOnly) {
              yaml += `
          readOnly: true`
            }
          }
        }
      }
    }
  }

  // 添加主容器配置
  yaml += `
      containers:
      - name: ${form.containerName}
        image: ${form.image}`

  // 添加端口配置
  if (form.ports.length > 0) {
    yaml += `
        ports:
${form.ports.map(p => `        - containerPort: ${p.containerPort}
          protocol: ${p.protocol}`).join('\n')}`
  }

  // 添加环境变量
  const validEnvs = form.env.filter(env => env.name && (env.type === 'value' ? env.value : env.type === 'configMapKeyRef' ? env.valueFrom.configMapKeyRef.name && env.valueFrom.configMapKeyRef.key : env.valueFrom.secretKeyRef.name && env.valueFrom.secretKeyRef.key))
  if (validEnvs.length > 0) {
    yaml += `
        env:`
    for (const env of validEnvs) {
      yaml += `
        - name: ${env.name}`
      if (env.type === 'value') {
        yaml += `
          value: ${env.value}`
      } else if (env.type === 'configMapKeyRef') {
        yaml += `
          valueFrom:
            configMapKeyRef:
              name: ${env.valueFrom.configMapKeyRef.name}
              key: ${env.valueFrom.configMapKeyRef.key}`
      } else if (env.type === 'secretKeyRef') {
        yaml += `
          valueFrom:
            secretKeyRef:
              name: ${env.valueFrom.secretKeyRef.name}
              key: ${env.valueFrom.secretKeyRef.key}`
      }
    }
  }

  // 添加生命周期配置
  const hasLifecycle = form.lifecycle.startupProbe.type || form.lifecycle.livenessProbe.type || form.lifecycle.readinessProbe.type
  if (hasLifecycle) {
    yaml += `
        lifecycle:`
    // 添加启动探针
    if (form.lifecycle.startupProbe.type) {
      yaml += `
        startupProbe:`
      if (form.lifecycle.startupProbe.type === 'httpGet') {
        yaml += `
          httpGet:`
        yaml += `
            path: ${form.lifecycle.startupProbe.httpGet.path}`
        yaml += `
            port: ${form.lifecycle.startupProbe.httpGet.port}`
        yaml += `
            scheme: ${form.lifecycle.startupProbe.httpGet.scheme}`
      } else if (form.lifecycle.startupProbe.type === 'tcpSocket') {
        yaml += `
          tcpSocket:`
        yaml += `
            port: ${form.lifecycle.startupProbe.tcpSocket.port}`
      } else if (form.lifecycle.startupProbe.type === 'exec') {
        yaml += `
          exec:`
        if (form.lifecycle.startupProbe.exec.command[0]) {
          yaml += `
            command:
            - ${form.lifecycle.startupProbe.exec.command[0]}`
        }
      }
      yaml += `
          initialDelaySeconds: ${form.lifecycle.startupProbe.initialDelaySeconds}`
      yaml += `
          periodSeconds: ${form.lifecycle.startupProbe.periodSeconds}`
      yaml += `
          timeoutSeconds: ${form.lifecycle.startupProbe.timeoutSeconds}`
      yaml += `
          successThreshold: ${form.lifecycle.startupProbe.successThreshold}`
      yaml += `
          failureThreshold: ${form.lifecycle.startupProbe.failureThreshold}`
    }
    // 添加存活探针
    if (form.lifecycle.livenessProbe.type) {
      yaml += `
        livenessProbe:`
      if (form.lifecycle.livenessProbe.type === 'httpGet') {
        yaml += `
          httpGet:`
        yaml += `
            path: ${form.lifecycle.livenessProbe.httpGet.path}`
        yaml += `
            port: ${form.lifecycle.livenessProbe.httpGet.port}`
        yaml += `
            scheme: ${form.lifecycle.livenessProbe.httpGet.scheme}`
      } else if (form.lifecycle.livenessProbe.type === 'tcpSocket') {
        yaml += `
          tcpSocket:`
        yaml += `
            port: ${form.lifecycle.livenessProbe.tcpSocket.port}`
      } else if (form.lifecycle.livenessProbe.type === 'exec') {
        yaml += `
          exec:`
        if (form.lifecycle.livenessProbe.exec.command[0]) {
          yaml += `
            command:
            - ${form.lifecycle.livenessProbe.exec.command[0]}`
        }
      }
      yaml += `
          initialDelaySeconds: ${form.lifecycle.livenessProbe.initialDelaySeconds}`
      yaml += `
          periodSeconds: ${form.lifecycle.livenessProbe.periodSeconds}`
      yaml += `
          timeoutSeconds: ${form.lifecycle.livenessProbe.timeoutSeconds}`
      yaml += `
          successThreshold: ${form.lifecycle.livenessProbe.successThreshold}`
      yaml += `
          failureThreshold: ${form.lifecycle.livenessProbe.failureThreshold}`
    }
    // 添加就绪探针
    if (form.lifecycle.readinessProbe.type) {
      yaml += `
        readinessProbe:`
      if (form.lifecycle.readinessProbe.type === 'httpGet') {
        yaml += `
          httpGet:`
        yaml += `
            path: ${form.lifecycle.readinessProbe.httpGet.path}`
        yaml += `
            port: ${form.lifecycle.readinessProbe.httpGet.port}`
        yaml += `
            scheme: ${form.lifecycle.readinessProbe.httpGet.scheme}`
      } else if (form.lifecycle.readinessProbe.type === 'tcpSocket') {
        yaml += `
          tcpSocket:`
        yaml += `
            port: ${form.lifecycle.readinessProbe.tcpSocket.port}`
      } else if (form.lifecycle.readinessProbe.type === 'exec') {
        yaml += `
          exec:`
        if (form.lifecycle.readinessProbe.exec.command[0]) {
          yaml += `
            command:
            - ${form.lifecycle.readinessProbe.exec.command[0]}`
        }
      }
      yaml += `
          initialDelaySeconds: ${form.lifecycle.readinessProbe.initialDelaySeconds}`
      yaml += `
          periodSeconds: ${form.lifecycle.readinessProbe.periodSeconds}`
      yaml += `
          timeoutSeconds: ${form.lifecycle.readinessProbe.timeoutSeconds}`
      yaml += `
          successThreshold: ${form.lifecycle.readinessProbe.successThreshold}`
      yaml += `
          failureThreshold: ${form.lifecycle.readinessProbe.failureThreshold}`
    }
  }

  // 添加资源配置
  yaml += `
        resources:
          requests:
            cpu: ${form.resources.requests.cpu}
            memory: ${form.resources.requests.memory}
          limits:
            cpu: ${form.resources.limits.cpu}
            memory: ${form.resources.limits.memory}`

  // 添加卷挂载
  const validVolumeMounts = form.volumeMounts.filter(mount => mount.name && mount.mountPath)
  if (validVolumeMounts.length > 0) {
    yaml += `
        volumeMounts:`
    for (const mount of validVolumeMounts) {
      yaml += `
        - name: ${mount.name}`
      yaml += `
          mountPath: ${mount.mountPath}`
      if (mount.subPath) {
        yaml += `
          subPath: ${mount.subPath}`
      }
      if (mount.readOnly) {
        yaml += `
          readOnly: true`
      }
    }
  }

  // 添加优雅终止时间
  if (form.terminationGracePeriodSeconds !== 30) {
    yaml += `
      terminationGracePeriodSeconds: ${form.terminationGracePeriodSeconds}`
  }

  yamlContent.value = yaml
}

// 预检YAML
const precheckYaml = async () => {
  if (!selectedCluster.value || !yamlContent.value.trim()) {
    ElMessage.warning('请选择集群并填写YAML内容')
    return
  }

  validationLoading.value = true
  validationResult.value = null

  try {
    const response = await axios.post(`/clusters/${selectedCluster.value}/deployments/precheck`, {
      yaml: yamlContent.value,
      namespace: deploymentForm.value.namespace
    })
    validationResult.value = `✅ 恭喜🎉：${response.data.message || 'YAML 定义符合规范，可以正常创建'}`
    ElMessage.success('YAML验证通过')
  } catch (error: any) {
    validationResult.value = `❌ 很遗憾😮‍💨：${error.response?.data?.error || error.message || '未知错误'}`
    ElMessage.error('YAML验证失败')
  } finally {
    validationLoading.value = false
  }
}

// 创建Deployment
const createDeployment = async () => {
  if (!selectedCluster.value) {
    ElMessage.warning('请选择一个集群')
    return
  }

  createLoading.value = true

  try {
    let deploymentData: any

    if (createMode.value === 'form') {
      // 表单模式
      // 先生成YAML，然后发送完整的YAML到后端
      generateYamlFromForm()
      deploymentData = {
        yaml: yamlContent.value,
        namespace: deploymentForm.value.namespace
      }
    } else {
      // YAML模式
      deploymentData = {
        yaml: yamlContent.value,
        namespace: deploymentForm.value.namespace
      }
    }

    // 调用创建API
    await axios.post(`/clusters/${selectedCluster.value}/deployments`, deploymentData)
    
    ElMessage.success('Deployment创建成功')
    createDeploymentVisible.value = false
    // 刷新Deployment列表
    if (activeWorkloadTab.value === 'deployment') {
      loadDeployments()
    }
  } catch (error: any) {
    console.error('创建Deployment失败:', error)
    ElMessage.error(`创建Deployment失败：${error.response?.data?.error || error.message || '未知错误'}`)
  } finally {
    createLoading.value = false
  }
}
</script>

<style scoped>
.deployments-view {
  width: 100%;
  padding: 0;
}

.page-header {
  margin-bottom: 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-content h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #333;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.deployments-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0;
  padding-bottom: 0;
}

/* 标签页样式 */
.tab-label {
  display: flex;
  align-items: center;
  gap: 8px;
}

.tab-label .el-icon {
  font-size: 16px;
}

/* 调整标签页间距和样式 */
:deep(.el-tabs__header) {
  margin-bottom: 0;
  padding-bottom: 0;
}

:deep(.el-tabs__nav-wrap) {
  padding-bottom: 0;
}

:deep(.el-tabs__item) {
  font-size: 14px;
  font-weight: 500;
  padding: 12px 20px;
}

:deep(.el-tabs__item.is-active) {
  color: var(--el-color-primary);
  font-weight: 600;
}

:deep(.el-tabs__active-bar) {
  height: 3px;
}

/* 列标题样式 */
.column-header {
  display: flex;
  align-items: center;
  gap: 5px;
}

/* 信息图标样式 */
.info-icon {
  font-size: 14px;
  color: #909399;
  cursor: help;
  transition: color 0.3s ease;
}

.info-icon:hover {
  color: #409eff;
}

/* 调整表格样式 */
:deep(.el-table) {
  margin-top: 0;
}

.deployment-details {
  padding: 10px 0;
}

.scale-deployment {
  padding: 10px 0;
}

.no-cluster-selected {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 300px;
  background-color: #fafafa;
  border-radius: 8px;
  margin-top: 20px;
}
</style>