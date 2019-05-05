@Library('automation@v3.0.1')

def label = "test_gotry_${UUID.randomUUID().toString()}"

// The normal go image since it contains make and we need it.
String image = 'golang:1.11'

// Each step the pipeline shall execute.
Map<String,Closure> steps = [
    "Test": {
        sh 'make coverage'
        publishHTML(target: [
            allowMissing: false,
            alwaysLinkToLastBuild: false,
            keepAll: false,
            reportDir: 'coverage',
            reportFiles: 'index.html',
            reportName: 'Coverage Report'
        ])
    }
]

genericAutomation.generic(label, image, steps)
