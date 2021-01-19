package com.example.gdsudao.activity

import android.app.DatePickerDialog
import android.content.Intent
import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.util.Log
import android.view.Menu
import android.view.MenuItem
import android.widget.Toast
import com.example.gdsudao.R
import kotlinx.android.synthetic.main.activity_detalhes.*
import java.text.SimpleDateFormat
import java.util.*
import kotlin.math.roundToInt

class DetalhesActivity : AppCompatActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_detalhes)

        setSupportActionBar(findViewById(R.id.toolbarDetalhes))

        // Get the areas
        var sp = com.example.gdsudao.utils.SharedPreferences()
        val bundle = intent.extras
        var areas = sp.RecuperarListaAreas(this)
        val areaIndex = bundle?.getInt("item")

        if(areaIndex != null){
            // Toast.makeText(this, "Area: ${areaIndex}", Toast.LENGTH_SHORT).show(); FUNCIONANDO
            Log.println(Log.DEBUG, "area", "Area: ${areas[areaIndex].toString()}")

            val ST_PRIRO_CORTE = 358.0f
            val ST_OUTROS_CORTES = 281.0f

            tvNomeArea.text = areas[areaIndex].nome

            if (!areas[areaIndex].proxcorte.isNullOrEmpty()){
                tvDataEstimada.text = "${areas[areaIndex].proxcorte.substring(8, 10)}/${areas[areaIndex].proxcorte.substring(5, 7)}/" +
                        "${areas[areaIndex].proxcorte.substring(0, 4)}"
            }else{
                tvDataEstimada.text = "ERRO"
                Toast.makeText(this, "Essa estação não contém dados inválidos!", Toast.LENGTH_SHORT).show()
            }


            if (!areas[areaIndex].dataCorte.isNullOrEmpty()){
                tvDataUltimoPastejo.text = "${areas[areaIndex].dataCorte.substring(8, 10)}/${areas[areaIndex].dataCorte.substring(5, 7)}/" +
                        "${areas[areaIndex].dataCorte.substring(0, 4)}"
            }else{
                tvDataUltimoPastejo.text = "ERRO"
            }

            /*
            if (areas[areaIndex].diario.length > 5){
                tvDiarios.text = areas[areaIndex].diario.substring(0, 4) + "%" //.substring(5) + "%"
            }else{
                tvDiarios.text = areas[areaIndex].diario + "%"


            if (areas[areaIndex].previsao.length > 5){
                tvPrevisoes.text = areas[areaIndex].previsao.substring(0, 4) + "%"  //.substring(5) + "%"
            }else{
                tvPrevisoes.text = areas[areaIndex].previsao + "%"  //.substring(5) + "%"
            }

            if (areas[areaIndex].normal.length > 5){
                tvNormais.text = areas[areaIndex].normal.substring(0, 4) + "%"  //.substring(5) + "%"
            }else{
                tvNormais.text = areas[areaIndex].normal + "%"  //.substring(5) + "%"
            }

            if (areas[areaIndex].st.length > 5){
                tvGdAcumulado.text = areas[areaIndex].st.substring(0, 5)  //.substring(5) + "%"
            }else{
                tvGdAcumulado.text = areas[areaIndex].st //.substring(5) + "%"
            }

            }
             */
            var dia = areas[areaIndex].diario.toFloat().roundToInt()
            if (dia != null){
                tvDiarios.text = dia.toString() + "%"
            }

            var pre = areas[areaIndex].previsao.toFloat().roundToInt()
            if (pre != null){
                tvPrevisoes.text = pre.toString() + "%"
            }

            var nor = areas[areaIndex].normal.toFloat().roundToInt()
            if (nor != null){
                tvNormais.text = nor.toString() + "%"
            }

            var st = areas[areaIndex].st.toFloat().roundToInt()
            if (st != null){
                tvGdAcumulado.text = st.toString() + " gd"
            }


            tvEstacao.text  = areas[areaIndex].codigoEstacao
            tvNumeroCortes.text = areas[areaIndex].numeroCorte

            if (!areas[areaIndex].st.isNullOrEmpty() && !areas[areaIndex].st.isNullOrEmpty()){
                var progresso: Float
                if (areas[areaIndex].numeroCorte.toInt() < 1) {
                    progresso = (areas[areaIndex].st.toFloat() / ST_PRIRO_CORTE) * 100
                }else {
                    progresso =  (areas[areaIndex].st.toFloat() / ST_OUTROS_CORTES ) * 100
                }


                progressBar.progress = progresso.toInt()

                /*
                if (progresso.toString().length > 5){
                    tvProgresso.text = progresso.toString().substring(0, 4) + "%"  //.substring(5) + "%"
                }else{
                    tvProgresso.text = progresso.toString() + "%"  //.substring(5) + "%"
                }
                */

            }else{
                progressBar.progress = 0
                //tvProgresso.text = "0%"
            }
            /*



        tvNome?.text = area.nome
        tvProxCorte?.text = area.proxcorte
        tvNumeroCortes?.text = "" + area.numeroCorte




        if (!area.proxcorte.isNullOrEmpty()){
            tvProxCorte?.text = "${area.proxcorte.substring(8, 10)}/${area.proxcorte.substring(5, 7)}/${area.proxcorte.substring(0, 4)}"
            //tvProxCorte?.text = area.proxcorte.substring(0, 10).replace("-","/", false)
        }else{
            tvProxCorte?.text = "ERRO"
        }


        if (area.diario.length > 5){
            tvDiarios?.text = area.diario.substring(0, 4) + "%" //.substring(5) + "%"
        }else{
            tvDiarios?.text = area.diario + "%"
        }

        if (area.previsao.length > 5){
            tvPrevisoes?.text = area.previsao.substring(0, 4) + "%"  //.substring(5) + "%"
        }else{
            tvPrevisoes?.text = area.previsao + "%"  //.substring(5) + "%"
        }

        if (area.normal.length > 5){
            tvNormais?.text = area.normal.substring(0, 4) + "%"  //.substring(5) + "%"
        }else{
            tvNormais?.text = area.normal + "%"  //.substring(5) + "%"
        }

        // tvProgresso?.text = "Data Plantio: " + "|d: " + area.diario + "|p: " + area.previsao +  "|n: " + area.normal




             */


        }else{
            Toast.makeText(this, "Erro ao apresentar os detalhes.", Toast.LENGTH_SHORT).show()
        }

    }

    override fun onCreateOptionsMenu(menu: Menu?): Boolean {
        menuInflater.inflate(R.menu.detalhes_menu, menu)
        return super.onCreateOptionsMenu(menu)
    }

    override fun onOptionsItemSelected(item: MenuItem): Boolean {
        var id = item.title

        // ADICIONAR NOVO PASTEJO
        if (id == "Adicionar pastejo"){


               var cal = Calendar.getInstance()

                val dateSetListener = DatePickerDialog.OnDateSetListener { view, year, monthOfYear, dayOfMonth ->
                cal.set(Calendar.YEAR, year)
                cal.set(Calendar.MONTH, monthOfYear)
                cal.set(Calendar.DAY_OF_MONTH, dayOfMonth)

                val dataFmt = SimpleDateFormat("yyyy-MM-dd", Locale.US).format(cal.time)



                val sp = com.example.gdsudao.utils.SharedPreferences()
                val bundle = intent.extras
                var areas = sp.RecuperarListaAreas(this)
                val areaIndex = bundle?.getInt("item")


                if (areaIndex != null){
                    //Toast.makeText(this, "${dataFmt} == ${areas[areaIndex].dataCorte}", Toast.LENGTH_SHORT).show()

                    val sdf = SimpleDateFormat("yyyy-MM-dd", Locale.US)
                    var ultimoPastejo = sdf.parse(areas[areaIndex].dataCorte)
                    var novoPastejo = sdf.parse(dataFmt)

                    if ( (areas[areaIndex].numeroCorte.toInt() < 1 && areas[areaIndex].st.toFloat() >= 200.00)
                        ||(areas[areaIndex].numeroCorte.toInt() >= 1 && areas[areaIndex].st.toFloat() >= 100.00)){


                        if (novoPastejo.after(ultimoPastejo)){
                            areas[areaIndex].dataCorte = dataFmt
                            val nc = areas[areaIndex].numeroCorte.toInt() + 1
                            areas[areaIndex].numeroCorte = nc.toString()

                            sp.AtualizarAreaLocal(this, areas[areaIndex], areaIndex)

                            var intent = Intent(this, MenuActivity::class.java)
                            startActivity(intent)
                        }else{
                            Toast.makeText(this, "O novo pastejo deve acontecer depois do último cadastrado.", Toast.LENGTH_SHORT).show()
                        }
                    }else{
                        if(areas[areaIndex].numeroCorte.toInt() < 1){
                            Toast.makeText(this, "O novo pastejo não deve ocorrer com menos de 200 graus-dia.", Toast.LENGTH_SHORT).show()
                        }else{
                            Toast.makeText(this, "O novo pastejo não deve ocorrer com menos de 100 graus-dia.", Toast.LENGTH_SHORT).show()
                        }
                    }

                }else{
                    Toast.makeText(this, "Erro ao cadastrar novo corte", Toast.LENGTH_SHORT).show()
                }
            }

            DatePickerDialog(this, dateSetListener,
                    cal.get(Calendar.YEAR),
                    cal.get(Calendar.MONTH),
                    cal.get(Calendar.DAY_OF_MONTH)).show()
        }

        // EXCLUIR AREA ATUAL
        if (id == "Excluir área"){
            val sp = com.example.gdsudao.utils.SharedPreferences()
            val bundle = intent.extras
            var areas = sp.RecuperarListaAreas(this)
            val areaIndex = bundle?.getInt("item")
            if(areaIndex != null){
                sp.RemoverAreaLista(this, areas[areaIndex])
                var intent = Intent(this, MenuActivity::class.java)
                startActivity(intent)
            }else{
                Toast.makeText(this, "Erro ao excluir area.", Toast.LENGTH_SHORT).show()
            }
        }

        return super.onOptionsItemSelected(item)
    }
}