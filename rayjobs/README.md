# Submit Ray Jobs

Install python: `brew install python`
Add python to zshrc: `echo "alias python=/usr/bin/python3" >> ~/.zshrc`
Set up virtual environment: `python -m venv venv`
Active virtual environment: `source venv/bin/activate`
Install ray: `pip install -U "ray[default]"`

Test if ray is installed: `ray --version`

## Submit Ray Job

Port forward ray head service:
`kubectl port-forward service/raycluster-kuberay-head-svc 8265:8265`

Access Ray dashboard: `http://localhost:8265`

Submit a test job:
`ray job submit --address http://localhost:8265 --working-dir ${PWD}/script -- python script.py`