const core = require('@actions/core');
const github = require('@actions/github');

const token_github = core.getInput(`token_action`, { required: true })
const repos = core.getInput(`repos`, { required: true })
const assignees = core.getInput(`assignees`, { required: true })

const oc = github.getOctokit(token_github)

async function getNoAssigneeIssues(repo_c) {
    core.info(`-------------------Start to get issues----------------`)
    let issues_all = new Array()
    let issues = await oc.rest.issues.listForRepo(
        {
            owner: repo_c.owner,
            repo:repo_c.repo,
            state: "open",
            assignee: "none",
            per_page: 100,
            page:1
        }
    )
    issues_all.push(...issues.data)
    core.info(`Round 1 finished, the number of this round: ${issues.data.length}`)
    let count = 2
    while (issues.data.length == 100) {
        issues = await oc.rest.issues.listForRepo(
            {
                owner: repo_c.owner,
                repo:repo_c.repo,
                state: "open",
                assignee: "none",
                per_page: 100,
                page:count
            }
        )
        issues_all.push(...issues.data)
        core.info(`Round ${count} finished, the number of this round: ${issues.data.length}`)
        count++
    }
    return issues_all
}

async function addAssignee(repo_c,number,assginees) {
    let {status: status} = await oc.rest.issues.addAssignees(
        {
            ...repo_c,
            issue_number: number,
            assignees: assginees
        }
    )
    return status
}


async function run() {
    let arr_repos = repos.split(`,`)
    let arr_assignees = assignees.split(`,`)
    for (let i = 0; i < arr_repos.length; i++) {
        const repo = arr_repos[i].split(`/`);
        let issues = await getNoAssigneeIssues(
            {
                owner: repo[0],
                repo:repo[1]
            }
        )
        core.info(`the number of issues which not assignee is: ${issues.length}`)
        for (const issue of issues) {
            if (issue.pull_request !== undefined) {
                core.info(`This issue ${issue.number} is pull request, so skip it`)
                continue
            }
            core.info(`Add assignees ${assignees} to issue ${issue.number}...`)
            let count = 0
            while (await addAssignee(repo,issue.number,arr_assignees) != 201 && count < 10) {
                core.info(`>>> failed, will try again... ${count+1}`)
            }
            if (count >= 10) {
                core.info(`>>>>>try 10 times, skip this issue ${issue.number}`)
            }else {
                core.info(`>>>>>success`)
            }
            core.info(``)
        }
    }
}

run()