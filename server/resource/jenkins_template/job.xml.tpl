<?xml version='1.1' encoding='UTF-8'?>
<!--
workflow-job plugin version
workflow-job@1316.vd2290d3341a_f Requires Jenkins 2.361.4
https://plugins.jenkins.io/workflow-job/releases/
-->
<flow-definition plugin="workflow-job@1316.vd2290d3341a_f">
    <actions/>
    <description>{{ .Description }}</description>
    <keepDependencies>false</keepDependencies>
    <properties>
        <jenkins.model.BuildDiscarderProperty>
            <strategy class="hudson.tasks.LogRotator">
                <daysToKeep>365</daysToKeep>
                <numToKeep>-1</numToKeep>
                <artifactDaysToKeep>-1</artifactDaysToKeep>
                <artifactNumToKeep>-1</artifactNumToKeep>
            </strategy>
        </jenkins.model.BuildDiscarderProperty>
        <hudson.model.ParametersDefinitionProperty>
            <parameterDefinitions>
        {{ range .BuildParams }}

                <hudson.model.StringParameterDefinition>
                    <name>{{ . }}</name>
                    <trim>false</trim>
                </hudson.model.StringParameterDefinition>
        {{ end }}

            </parameterDefinitions>
        </hudson.model.ParametersDefinitionProperty>
    </properties>
    <!--
  workflow-cps plugin version
  workflow-cps@3705.va_6a_c2775a_c17 Requires Jenkins 2.387.3
  https://plugins.jenkins.io/workflow-cps/
  -->
    <definition class="org.jenkinsci.plugins.workflow.cps.CpsScmFlowDefinition" plugin="workflow-cps@3705.va_6a_c2775a_c17">
        <!--
    git plugin version
    git@5.2.0 Requires Jenkins 2.387.3
    https://plugins.jenkins.io/git/
    -->
        <scm class="hudson.plugins.git.GitSCM" plugin="git@5.2.0">
            <configVersion>2</configVersion>
            <userRemoteConfigs>
                <hudson.plugins.git.UserRemoteConfig>
                    <url>{{ .GitRepo }}</url>
                    <credentialsId>{{ .GitCredentialId }}</credentialsId>
                </hudson.plugins.git.UserRemoteConfig>
            </userRemoteConfigs>
            <branches>
                <hudson.plugins.git.BranchSpec>
                    <name>*/${branch}</name>
                </hudson.plugins.git.BranchSpec>
            </branches>
            <doGenerateSubmoduleConfigurations>false</doGenerateSubmoduleConfigurations>
            <submoduleCfg class="empty-list"/>
            <extensions/>
        </scm>
        <scriptPath>{{ .Jenkinsfile }}</scriptPath>
        <lightweight>true</lightweight>
    </definition>
    <triggers/>
    <quietPeriod>1</quietPeriod>
    <disabled>false</disabled>
</flow-definition>