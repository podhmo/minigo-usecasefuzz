#!/usr/bin/env python3
import argparse, concurrent.futures, json, os, pathlib, subprocess
root=pathlib.Path(__file__).resolve().parent
parser=argparse.ArgumentParser()
parser.add_argument('--minigo',required=True)
parser.add_argument('cases',nargs='*')
args=parser.parse_args()
cases=args.cases or sorted(p.name for p in root.iterdir() if (p/'main.go').exists() and p.name not in {'session-policy','type-shadow','float-width'})
env=dict(os.environ,GOCACHE='/private/tmp/minigo-review-cache')
def run(p,cmd):
 try:
  r=subprocess.run(cmd,cwd=p,env=env,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,timeout=20);return r.returncode,r.stdout.decode(errors='replace')
 except subprocess.TimeoutExpired:return 124,'TIMEOUT\n'
def compare(name):
 p=root/name;g=run(p,['go','run','.']);m=run(p,[str(pathlib.Path(args.minigo).resolve()),'run','.']);(p/'go.stdout').write_text(g[1]);(p/'minigo.stdout').write_text(m[1]);return {'name':name,'go_rc':g[0],'minigo_rc':m[0],'go':g[1].strip(),'minigo':m[1].strip()}
with concurrent.futures.ThreadPoolExecutor(max_workers=4) as ex:results=list(ex.map(compare,cases))
(root/'verified-results.json').write_text(json.dumps(results,indent=2))
for r in results:
 if r['go']!=r['minigo'] or r['go_rc']!=r['minigo_rc']:print(r['name']+': Go='+repr(r['go'])+'; minigo='+repr(r['minigo']))
